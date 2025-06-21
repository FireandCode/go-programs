//////////////////////////////////////////////////////////////////////
//
// Given is some code to cache key-value pairs from a database into
// the main memory (to reduce access time). Note that golang's map are
// not entirely thread safe. Multiple readers are fine, but multiple
// writers are not. Change the code to make this thread safe.
//

package main

import (
	"container/list"
	"sync"
	"testing"

	"golang.org/x/sync/singleflight"
)

var sngfl singleflight.Group 

// CacheSize determines how big the cache can grow
const CacheSize = 100

// KeyStoreCacheLoader is an interface for the KeyStoreCache
type KeyStoreCacheLoader interface {
	// Load implements a function where the cache should gets it's content from
	Load(string) string
}

type page struct {
	Key   string
	Value string
}

// KeyStoreCache is a LRU cache for string key-value pairs
type KeyStoreCache struct {
	cache map[string]*list.Element
	pages list.List
	load  func(string) string
	mul sync.RWMutex
}

// New creates a new KeyStoreCache
func New(load KeyStoreCacheLoader) *KeyStoreCache {
	return &KeyStoreCache{
		load:  load.Load,
		cache: make(map[string]*list.Element),
	}
}

// Get gets the key from cache, loads it from the source if needed
func (k *KeyStoreCache) Get(key string) string {

	k.mul.RLock()
	if e, ok := k.cache[key]; ok {
		val := e.Value.(page).Value
		k.mul.RUnlock()

		k.mul.Lock()
		k.pages.MoveToFront(e)
		k.mul.Unlock()

		return val 
	}
	k.mul.RUnlock()
	

	// Miss - load from database and save it in cache
	val, _ , _ := sngfl.Do(key, func() (interface{}, error) {
		val := k.load(key)
		
		return val, nil
	})
	valSt := val.(string)
	p := page{Key: key, Value: valSt}

	// if cache is full remove the least used item

	k.mul.Lock()
	defer k.mul.Unlock()

		// Double-check if another thread already inserted it
	if e, ok := k.cache[key]; ok {
		k.pages.MoveToFront(e)
		return e.Value.(page).Value
	}

	if len(k.cache) >= CacheSize {
		end := k.pages.Back()
		if end != nil {
			evicted := end.Value.(page)
			delete(k.cache, evicted.Key)
			k.pages.Remove(end)
		}
	}


	elem := k.pages.PushFront(p)
	k.cache[key] = elem
	return valSt
}

// Loader implements KeyStoreLoader
type Loader struct {
	DB *MockDB
}

// Load gets the data from the database
func (l *Loader) Load(key string) string {
	val, err := l.DB.Get(key)
	if err != nil {
		panic(err)
	}

	return val
}

func run(t *testing.T) (*KeyStoreCache, *MockDB) {
	loader := Loader{
		DB: GetMockDB(),
	}
	cache := New(&loader)

	RunMockServer(cache, t)

	return cache, loader.DB
}

func main() {
	run(nil)
}
