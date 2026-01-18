package main

import (
	"fmt"
)

/*
1. go through the golang fundamentals
2. go throught the design patterns
3. solve the problems using the design patterns

golang fundamentals

basic level
data structures

medium level
go routines, mutexes, channels, sync library, time library

advance level
concurrent operations, atomic operations, semaphores, memory calculation, time calculation, first flight, routers, frameworks, memory management, garbage collector

design patterns
on my own i will go through the design patterns

golang problems



basic level ->
data structures
1. primary/basic
int ,string, bool,float,double, enum, channel
2. aggregated
array, slice, list, map, priority_queue, segment tree, fenwick tree,
3. key words
function, struct, pointers, reference by memory/value, special functions

medium level ->
go routines -> go keyword
mutexes -> sync.Mutex, sync.RMutex
channels -> buffer, unbuffered
string ->
conversion library ->
sync ->
time ->

concurrent operations, atomic operations, semaphores, memory calculation, time calculation, single flight, routers, frameworks, memory management, garbage collector
advance level
concurrent operations -> go keyword to make go routines to make a program/function concurrent.
waitgroups :-
atomic operations -> sync.atomic to make a particular operations atomic used in concurrent settings where we need
to update the value of a particular data for example counter and using mutexes is overdoing it
semaphores -> semaphores are atomic counters which have certain threshold. which can be used to used in situation
where we want to run
memory calculations -> runtime library,
time calculation -> props
single flight -> we have cache and DB, 100/requests per second, cache 5 sec expire.
if cache expires -> per request it takes 1 second to get the value from DB.
having this 100 requests going to DB -> single flight -> mutex/semaphore will be taken request -> hit the DB -> update the cache
routers :- router library
frameworks :- gin, gorm,
garbage collector/memory management :- algorithm -> tricolor algorithm -> white, gray, black
each time a garbage collector -> if memory is directly accessed -> black, if it is referenced -> gray
if it not at all reacheable -> then it will remain white
-> garbage collector will remove all the white memory references.

->


basic level ->
data structures
1. primary/basic
int ,string, bool,float,double, enum, channel
2. aggregated
array, slice, list, map, priority_queue, segment tree, fenwick tree,
3. key words
function,type , interface{}, struct, pointers, reference by memory/value,
special functions
main – Program entry point, runs after all init() functions.

init – Auto-executed setup function, runs before main, cannot be called manually.

panic – Aborts normal execution and starts stack unwinding.

recover – Stops a panic and regains control, works only inside defer.

defer – Schedules a function to run after the surrounding function returns.

len – Returns the length of a data structure.

cap – Returns the capacity of slices, arrays, or channels.

make – Initializes slices, maps, and channels.

new – Allocates zeroed memory and returns a pointer.

append – Adds elements to a slice, reallocating if needed.

copy – Copies elements between slices.

delete – Removes a key from a map.

close – Closes a channel to signal no more values.

complex – Constructs a complex number.

real – Extracts the real part of a complex number.

imag – Extracts the imaginary part of a complex number.

*/


func init()  {
	fmt.Println("Hello world from the starting point")	
}

type myMap map[string]int

type IPayment interface{
	PaymentProcess() bool 
	SetX(x int)
}

type Payment struct {
	x int
	y *int
}

func(p *Payment) SetX(x int) {
	p.x = x
}

func (p *Payment) PaymentProcess() bool {
	p.x += 1
	fmt.Println(p.x)
	fmt.Println(p.y)
	fmt.Println(*p.y)
	*p.y = 12
	return true
}

func ReferenceBy(x int, y *int, pay *Payment) {
	x = 10
	*y = 122
	pay.SetX(x)
}

func main() {

	/*
	medium level ->
go routines -> go keyword
mutexes -> sync.Mutex, sync.RMutex
channels -> buffer, unbuffered
string ->
conversion library ->
sync ->
time ->
	
	*/
// 	xy := 11
// 	yx := 22
// 	 xx := 10
//  payment := Payment{
// 	x:0,
// y: &xx}

// xMMake :=int(22)
// xMNew := new(int)
// xMNew = &xy 
// xMMake = 12 

// fmt.Println(*xMNew, xMMake)
//  payment.PaymentProcess()
//  var pay IPayment  = &payment

// ReferenceBy(xy, &yx, &payment)
// fmt.Println(xy, yx, payment.x)



// pay.PaymentProcess()

//  var sliA2 = []int{23,22, 55, 5232, 23 }
//  fmt.Println(sliA2)
//  sliA2 = sliA2[1:3]
//  fmt.Println(sliA2)


	// fmt.Println("HELLO WORLD")
	// var arr = [3]int{11, 11 ,22 }
	// var sliA []int
	// var sliA2 = []int{23,22, 55, 5232, 23 }
	// sliA = append(sliA, 11,123, 11, 1323, 5, 55 )
	// fmt.Println(arr, sliA, sliA2)

	// xMap := make(myMap)
	// xMap["123"] = 123
	// xMap["323"] = 323

	// for k,v := range xMap {
	// 	fmt.Println(k,v)
	// }
// 	func PrintValue(x chan int) {
// 	x <- 111123
// }

	// slices.SortStableFunc(sliA2, func (x int ,y int) int  {
	// 	if x > y {
	// 		return -1
	// 	}
	// 	return 0
	// })
	// sort.Slice(sliA, func(i, j int) bool {
	// 	return sliA[i]< sliA[j]
	// })
	// fmt.Println(sort.Find(len(sliA), func(i int) int {
	// 	if sliA[i] > 123 {
	// 		 return -1
	// 	}
	// 	if sliA[i] < 123 {
	// 		return 1
	// 	}
	// 	return 0
	// }))
	// fmt.Println(sliA2, sliA)
	// l := list.New()
	// l.PushBack(11)
	// l.PushFront("string")
	// cur := l.Front()
	// for i:=0; i< l.Len(); i++ {
	// 	fmt.Println(cur.Value)
	// 	cur = cur.Next()
	// }






	// var x chan int 
	// y := make(chan int)
	// var x int 	
	
	// go PrintValue(y)

	// fmt.Println(<-y)
	// fmt.Println(x, "END WORLD")

}