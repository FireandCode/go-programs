package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

/*

EvictionPolicy - interface{}
  check()
  setValue()

  LRU - struct
  - map[Key][time.Time]
  - queue<time.Time> Key

  LFU - struct
  - map[Key]int
  - queue<time.Time> Key

- KeyTime struct
	- key Key
	- value - Time.Time

Key - struct
- value interface{}
- expiry time.Time

Value - struct
- value interface{}

for removing Expired Keys. we will have a go routine

go routines func
- removeExpiredKeys

DiskCache
- map[Key][Value]

Cache struct
- capacity
- map[Key]Value
- hits
- misses
- evictionCounts
- RWmutex


*/

type MultiCache struct {
	dc *DiskCache
	c *Cache
	mut *sync.Mutex
}

type DiskCache struct {
	store map[Key]Value
}


type Cache struct {
	capacity int
	store map[Key]Value
	epolicy EvictionPolicy
	hits int
	misses int
	evictionCounts int
}

type KeyTime struct {
	key Key
	value time.Time
}

type EvictionPolicy interface {
	check() Key
	setValue(Key)
	removeKey(Key)
	Initalize()
}

type LRU struct {
	m map[Key]*list.Element
	q *list.List
}

type Queue struct {
    items *list.List
}

func NewQueue() *Queue {
    return &Queue{items: list.New()}
}

// Enqueue (Push to the back)
func (q *Queue) Enqueue(value KeyTime) {
    q.items.PushBack(value)
}

// Dequeue (Pop from the front)
func (q *Queue) Dequeue() (KeyTime, bool) {
    if q.items.Len() == 0 {
        return KeyTime{}, false
    }
    front := q.items.Front()
    q.items.Remove(front)
    return front.Value.(KeyTime), true
}


type Key struct {
	value interface{}
	expiry time.Time
}

type Value struct {
	value interface{}
}

/*
  
*/
type LFU struct {
	capacity int
	store    map[Key]*list.Element
	freqMap  map[int]*list.List
	minFreq  int
	mut      *sync.Mutex
}

type Node struct {
	key   Key
	value Value
	freq  int
}

func (lfu *LFU) Initialize(capacity int) {
	lfu.capacity = capacity
	lfu.store = make(map[Key]*list.Element)
	lfu.freqMap = make(map[int]*list.List)
	lfu.minFreq = 0
	lfu.mut = &sync.Mutex{}
}

func (lfu *LFU) setValue(key Key) {
	lfu.mut.Lock()
	defer lfu.mut.Unlock()

	if elem, exists := lfu.store[key]; exists {
		node := elem.Value.(*Node)
		lfu.updateFrequency(node)
		return
	}

	if len(lfu.store) >= lfu.capacity {
		lfu.check()
	}

	node := &Node{key: key, freq: 1}
	if lfu.freqMap[1] == nil {
		lfu.freqMap[1] = list.New()
	}
	elem := lfu.freqMap[1].PushBack(node)
	lfu.store[key] = elem
	lfu.minFreq = 1
}

func (lfu *LFU) check() Key {
	list := lfu.freqMap[lfu.minFreq]
	if list == nil {
		return Key{}
	}
	elem := list.Front()
	if elem == nil {
		return Key{}
	}

	node := elem.Value.(*Node)
	delete(lfu.store, node.key)
	list.Remove(elem)

	if list.Len() == 0 {
		delete(lfu.freqMap, lfu.minFreq)
	}

	return node.key
}

func (lfu *LFU) removeKey(key Key) {

	lfu.mut.Lock()
	defer lfu.mut.Unlock()

	if elem, exists := lfu.store[key]; exists {
		node := elem.Value.(*Node)
		lfu.freqMap[node.freq].Remove(lfu.store[node.key])
		delete(lfu.store, node.key)
		return
	}

}

func (lfu *LFU) updateFrequency(node *Node) {
	oldFreq := node.freq
	node.freq++

	lfu.freqMap[oldFreq].Remove(lfu.store[node.key])
	if lfu.freqMap[oldFreq].Len() == 0 {
		delete(lfu.freqMap, oldFreq)
		if lfu.minFreq == oldFreq {
			lfu.minFreq++
		}
	}

	if lfu.freqMap[node.freq] == nil {
		lfu.freqMap[node.freq] = list.New()
	}
	lfu.store[node.key] = lfu.freqMap[node.freq].PushBack(node)
}



func (lr *LRU) setValue(key Key) {
	if elem, found := lr.m[key]; found {
		lr.q.MoveToFront(elem)
	} else {
		elem := lr.q.PushFront(key)
		lr.m[key] = elem
	}
}

func (lr *LRU) check() Key {
	elem := lr.q.Back()
	if elem != nil {
		key := elem.Value.(Key)
		lr.q.Remove(elem)
		delete(lr.m, key)
		return key
	}
	fmt.Println("No Key Present for Eviction Policy - LRU")
	return Key{}
}

func (lru *LRU) removeKey(key Key) {

	if elem, found := lru.m[key]; found {
		lru.q.Remove(elem)
		delete(lru.m, key)
	}
}

func (mu *MultiCache) ChangeCapacity(capacity int) {
	mu.mut.Lock()
	defer mu.mut.Unlock()
   for mu.c.capacity > capacity {
	evictedKey := mu.c.epolicy.check()
	mu.dc.store[evictedKey] = mu.c.store[evictedKey]
	mu.c.evictionCounts++

	delete(mu.c.store, evictedKey)
	mu.c.capacity--
}
}


func(mu *MultiCache) GetStatstics()  {
	
	fmt.Println("hits ", mu.c.hits)
	fmt.Println("misses ", mu.c.misses)
	fmt.Println("evictionCounts ", mu.c.evictionCounts)
}

func(mu *MultiCache) GetKey(key Key) Value {

	value, exists := mu.c.store[key]

	if exists {
		//check the evictionPolicy
		mu.c.hits++
		mu.c.epolicy.setValue(key) 
		return value
	}

	mu.c.misses++
	value, exists = mu.dc.store[key]
	if exists {
		return value
	}
	fmt.Println("Key requested is not present on cache and disk cache") 

	return Value{}
}

func(mu *MultiCache) removeKey(key Key)  {
	mu.mut.Lock()
	defer mu.mut.Unlock()
	delete(mu.c.store , key)
	mu.c.epolicy.setValue(key)
}


func(mu *MultiCache) SetKey(key Key, value Value) {
	//check the capacity
	mu.mut.Lock()
	defer mu.mut.Unlock()

	_, exists := mu.c.store[key]

	if !exists && mu.c.capacity == len(mu.c.store) {
		//check the evictionPolicy

		evictedKey := mu.c.epolicy.check()
		mu.dc.store[evictedKey] = mu.c.store[evictedKey]
		mu.c.evictionCounts++

		delete(mu.c.store, evictedKey)
		mu.c.store[key] = value
		mu.c.epolicy.setValue(key)

		return
	}

	mu.c.store[key] = value
	mu.c.epolicy.setValue(key) 
}

func NewDiskCache() *DiskCache {
	store := make(map[Key]Value)

	return &DiskCache{
		store: store,
	}
}

func(lr *LRU) Initalize() {
	lr.m = make(map[Key]*list.Element)
	lr.q = list.New()

}

func  NewCache(capacity int, ePolicy EvictionPolicy) *Cache  {
	store := make(map[Key]Value)
	return &Cache{
		capacity:  capacity,
		store: store,
		epolicy: ePolicy,
	}
}

func NewMultiCache(capacity int, ePolicy EvictionPolicy) *MultiCache  {
	
	dc := NewDiskCache()
	c := NewCache(capacity, ePolicy)
	mu := MultiCache{}

	mu.dc = dc
	mu.c = c
	mu.mut = &sync.Mutex{}
	return &mu
}


func main() {
	lru := LRU{}
	lru.Initalize()
	mu := NewMultiCache(3, &lru)

	go func(mu *MultiCache) {
		
		for {
			for k, v := range mu.c.store {
				if time.Now().After(k.expiry) {
					
					fmt.Println("Key ", k.value, " with value ", v, "is expired")
					mu.removeKey(k)
				}
			}
			time.Sleep(1* time.Second)
			
		}
	}(mu)
	mu.SetKey(Key{value: "1", expiry: time.Now().Add(3* time.Second)}, Value{1})
	mu.SetKey(Key{value: "2", expiry: time.Now().Add(1* time.Second)}, Value{1})
	mu.SetKey(Key{value: "3", expiry: time.Now().Add(2* time.Second)}, Value{1})
	mu.SetKey(Key{value: "4", expiry: time.Now().Add(3* time.Second)}, Value{1})

	mu.GetKey(Key{value: "2"})
	mu.GetStatstics()

		// Prevent main from exiting
		select {}
}