package main

import (
	"fmt"
	"sync"
	"time"
)

//barber problem

//shop, customer, barber, waiting room

//barber is done with each customer and a new customer comes

/*
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Customer represents a customer with a name
type Customer struct {
	name string
}

// Barber represents the barber
type Barber struct {
	name string
}

// Barber cuts hair of customers in the waiting room
func (b *Barber) cutHair(wr *WaitingRoom, stopCh chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case customer := <-wr.chair:
			fmt.Printf("%s is cutting hair for %s\n", b.name, customer.name)
			time.Sleep(cuttingTime)
			fmt.Printf("%s is done with %s\n", b.name, customer.name)
		case <-stopCh:
			fmt.Printf("%s is going home.\n", b.name)
			return
		}
	}
}

// WaitingRoom represents a waiting room with limited chairs
type WaitingRoom struct {
	chair chan Customer
}

// NewWaitingRoom initializes a waiting room with a fixed number of chairs
func NewWaitingRoom(size int) *WaitingRoom {
	return &WaitingRoom{
		chair: make(chan Customer, size),
	}
}

// AddCustomer tries to seat a customer in the waiting room
func (wr *WaitingRoom) addCustomer(customer Customer) bool {
	select {
	case wr.chair <- customer:
		fmt.Printf("%s is seated in the waiting room.\n", customer.name)
		return true
	default:
		fmt.Printf("%s found no empty chair and left.\n", customer.name)
		return false
	}
}

// Shop represents the barber shop
type Shop struct {
	barber      Barber
	waitingRoom *WaitingRoom
}

// NewShop initializes the shop with a barber and a waiting room
func NewShop(barberName string, waitingRoomSize int) *Shop {
	return &Shop{
		barber:      Barber{name: barberName},
		waitingRoom: NewWaitingRoom(waitingRoomSize),
	}
}

// Simulate simulates the barber shop with customer arrivals and the barber working
func (shop *Shop) simulate(customers []Customer) {
	var wg sync.WaitGroup
	stopCh := make(chan struct{})

	wg.Add(1)
	go shop.barber.cutHair(shop.waitingRoom, stopCh, &wg)

	for _, customer := range customers {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
		shop.waitingRoom.addCustomer(customer)
	}

	// Wait for remaining customers to be processed
	time.Sleep(time.Second)
	close(stopCh) // Signal the barber to stop
	wg.Wait()     // Wait for the barber to finish
	fmt.Println("The shop is now closed.")
}

const (
	cuttingTime = time.Second / 10
)

func main() {
	rand.Seed(time.Now().UnixNano())
	shop := NewShop("Barber Bob", 4)

	customers := []Customer{
		{name: "Rahul"},
		{name: "Rohit"},
		{name: "Chinmay"},
		{name: "Nitish"},
		{name: "Shubham"},
		{name: "Ram"},
	}

	shop.simulate(customers)
}

*/
//

// type Customer struct {
//     name string
// }

// type Barber struct {
//     name string
// }

// func (b *Barber) cutHair(wr *WaitingRoom, wg *sync.WaitGroup)  {

//     for i:=0; i<6; i++ {
//         name := <- wr.chair

//         fmt.Printf("%s is done\n", name.name)
//     }

//     wg.Done()
// }

// type WaitingRoom struct {
//     chair chan Customer
// }

// func InitializeWaitingRoom() WaitingRoom {
//     chairs := make(chan Customer, bufferSize)

//     return WaitingRoom{
//         chair: chairs,
//     }
// }
// func (wr *WaitingRoom) takeAChair(customer Customer)  {
//      if(len(wr.chair) < bufferSize ) {
//         wr.chair <- customer
//         fmt.Printf("%s is seated\n", customer.name)
//         return
//     }

//         fmt.Printf("%s has left\n", customer.name)

// }

// type Shop struct {
//     barber Barber
//     waitingRoom WaitingRoom
// }

// const (
//     cuttingTime = time.Second/10
//     bufferSize = 4
// )

// func main()  {
//     var wg sync.WaitGroup

//     wg.Add(1)
//     customers := []Customer{
//         {name: "Rahul"},
//         {name: "Rohit"},
//         {name: "Chinmay"},
//         {name: "Nitish"},
//         {name: "Shubham"},
//         {name: "Ram"},
//     }

//     shop := Shop{
//         barber: Barber{},
//         waitingRoom: InitializeWaitingRoom(),
//     }

//     for _, customer := range customers {
//            shop.waitingRoom.takeAChair(customer)
//     }

//     go shop.barber.cutHair(&shop.waitingRoom, &wg)

//     shop.waitingRoom.takeAChair(Customer{name:"somevalue"})
//     wg.Wait()
// }

//producer and consumer problem
// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// // Configuration constants
// const (
// 	bufferSize     = 5
// 	itemCount      = 10
// 	producerDelay  = time.Millisecond * 200
// 	consumerDelay  = time.Millisecond * 500
// )

// // Buffer struct to encapsulate buffer operations
// type Buffer struct {
// 	data chan int
// }

// // NewBuffer creates a new buffer with the given size
// func NewBuffer(size int) *Buffer {
// 	return &Buffer{data: make(chan int, size)}
// }

// // Produce adds an item to the buffer
// func (b *Buffer) Produce(item int) {
// 	b.data <- item
// 	fmt.Printf("Produced item: %d\n", item)
// }

// // Consume removes and returns an item from the buffer
// func (b *Buffer) Consume() int {
// 	item := <-b.data
// 	fmt.Printf("Consumed item: %d\n", item)
// 	return item
// }

// // Worker interface
// type Worker interface {
// 	Work(buffer *Buffer, wg *sync.WaitGroup)
// }

// // Producer struct
// type Producer struct {
// 	id int
// }

// // Work generates items and puts them into the buffer (implements Worker)
// func (p *Producer) Work(buffer *Buffer, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	fmt.Printf("Producer %d started\n", p.id)
// 	for i := 0; i < itemCount; i++ {
// 		time.Sleep(producerDelay)
// 		buffer.Produce(i)
// 	}
// 	fmt.Printf("Producer %d finished\n", p.id)
// }

// // Consumer struct
// type Consumer struct {
// 	id int
// }

// // Work retrieves items from the buffer (implements Worker)
// func (c *Consumer) Work(buffer *Buffer, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	fmt.Printf("Consumer %d started\n", c.id)
// 	for i := 0; i < itemCount; i++ {
// 		time.Sleep(consumerDelay)
// 		buffer.Consume()
// 	}
// 	fmt.Printf("Consumer %d finished\n", c.id)
// }

// func NewCompleteWorker(workerTypes []string) []Worker {
//     workers := make([]Worker, len(workerTypes))

//         for i, workerType := range workerTypes {
//             workers[i] = NewWorker(workerType, 1)
//         }

//     return workers
// }

// // NewWorker initializes and returns a Worker
// func NewWorker(workerType string, id int) Worker {
// 	switch workerType {
// 	case "producer":
// 		return &Producer{id: id}
// 	case "consumer":
// 		return &Consumer{id: id}
// 	default:
// 		panic("Invalid worker type")
// 	}
// }

// func main() {
// 	var wg sync.WaitGroup

// 	// Initialize buffer
// 	buffer := NewBuffer(bufferSize)

//     workerTypes := []string{"producer", "consumer"}
// 	// Create workers using NewWorker
// 	workers := NewCompleteWorker(workerTypes)

// 	// Add workers to WaitGroup and start their work
// 	wg.Add(len(workers))
// 	for _, worker := range workers {
// 		go worker.Work(buffer, &wg)
// 	}

// 	// Wait for all workers to complete
// 	wg.Wait()

// 	fmt.Println("All tasks are complete.")
// }

// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// //Consumer and Producer Problem

// /*
// producer, consumer - go routine
// buffer channel -
//     - producer will sent the value
//     - consumer will receive the value

// if channel is full
//     - consumer can consume it
//     - producer will go to sleep

// if channel is empty
//     - producer can sent the value
//     - consumer will go to sleep

// */

// const bufferSize = 5
// const itemSize = 10

// func producer(buffer chan int, wg *sync.WaitGroup)  {
//         defer wg.Done()
//         fmt.Println("Producer has started")
//         for i :=0; i<itemSize; i++ {
//             time.Sleep(time.Second/10)
//             buffer <- i
//             fmt.Printf("Producer has produced %d item\n", i)
//         }
// }

// func consumer(buffer chan int, wg *sync.WaitGroup) {
//     defer wg.Done()
//     fmt.Println("Consumer has started")

//     for i :=0; i< itemSize; i++ {
//         time.Sleep(time.Second/5)
//         val := <- buffer
//         fmt.Printf("Consumer has consumed %d item\n", val)
//     }

// }

// func main()  {
//     var wg *sync.WaitGroup = &sync.WaitGroup{}
//     buffer := make(chan int, bufferSize)

//     wg.Add(2)

//     go producer(buffer, wg)

//     go consumer(buffer, wg)

//     wg.Wait()

//     fmt.Println("Both producer and consumer are done")
// }

// import (
// 	"fmt"
// 	"hash/fnv"
// 	"math/rand"
// 	"sync"
// 	"time"
// )

// /*

// each worker is a go routine

// 10 computers
// for creating one computer i need 5 workers so 5 go routines

// for 10 times
//     wait.add(5)
//     for 5 workers
//      go routine initiate()

//     wait.wait()

// */

// func rGenerator(task string) *rand.Rand {
//  h := fnv.New64a()
//  h.Write([]byte(task))
//  rg := rand.New(rand.NewSource(int64(h.Sum64())))

//  return rg
// }

// func rSleep(rg *rand.Rand) {
//     time.Sleep(time.Second/2 + time.Duration(rg.Int63n(int64(time.Second/2))))
// }

// type Worker struct {
//     name string
//     task string
// }

// func(w *Worker) executeTask(isDone *sync.WaitGroup)  {
//     rg := rGenerator(w.task)

//     fmt.Printf("%s is working on %s\n", w.name, w.task)
//     rSleep(rg)

//     isDone.Done()
// }

// type Factory struct {
//     workers []Worker
// }

// func(f *Factory) Initialize(noOfWorkers int) {

//     workers := make([]Worker, noOfWorkers)
//     for i:=0; i< noOfWorkers; i++ {
//         workers[i].name = fmt.Sprintf("%d worker", i)
//         workers[i].task = fmt.Sprintf("%d task", i)
//     }

//     f.workers = workers
// }

// func (f * Factory) StartFactor(things int) {
//      for i :=0; i< things; i++ {
//         var isDone sync.WaitGroup
//         isDone.Add(len(f.workers))

//         for _, w := range f.workers {
//             go w.executeTask(&isDone)
//         }
//         fmt.Printf("%d thing is completed\n", i)

//         isDone.Wait()
//      }

//      fmt.Println("Factory has developed all the things")
// }

// func main()  {
//     var f Factory
//     f.Initialize(10)

//     f.StartFactor(7)
// }

//Philosophers Problem
// import (
// 	"fmt"
// 	"hash/fnv"
// 	"math/rand"
// 	"strconv"
// 	"sync"
// 	"time"
// )

// type Fork struct {
//     lc *sync.Mutex
// }

// type Philosopher struct {
//     Name string
//     DominantHand *Fork
//     OtherHand *Fork
// }

// type DiningTable struct {
//     Philosophers []Philosopher
//     isDone sync.WaitGroup
// }

// const times = 3
// const eating = time.Second/ 100
// const thinking = time.Second/100

// func rGenerator(name string) (*rand.Rand) {
//     h := fnv.New64a()
//     h.Write([]byte(name))
//     rg := rand.New(rand.NewSource(int64(h.Sum64())))

//     return rg
// }

// func rSleep(t time.Duration, rg *rand.Rand)  {

//     time.Sleep(t/2 + time.Duration(rg.Int63n(int64(t))))
// }

// func (p Philosopher) Eating() {
//     fmt.Printf("%s is Eating\n", p.Name)
// }

// func (p Philosopher) Thinking() {
//     fmt.Printf("%s is Thinking\n", p.Name)
// }

// func (p *Philosopher) Dine(rg *rand.Rand, isDone *sync.WaitGroup)  {
//     fmt.Println(p.Name, " has come to the table")

//     for i:=0; i<times ; i++ {

//         p.DominantHand.lc.Lock()
//         p.OtherHand.lc.Lock()
//         p.Eating()

//         rSleep(eating, rg)

//         p.DominantHand.lc.Unlock()
//         p.OtherHand.lc.Unlock()

//         p.Thinking()
//         rSleep(thinking, rg)
//     }

//    fmt.Println(p.Name , " has left the table")

//    isDone.Done()

// }

// func (dt *DiningTable) Initialize()  {
//     fork0 := &Fork{
//         &sync.Mutex{} ,
//     }
//     forkLeft := fork0
//     var philosophers []Philosopher
//     for i:=1; i<5; i++ {
//         fork1 :=  &Fork{
//             &sync.Mutex{} ,
//         }
//         philosopher := Philosopher{}
//         philosopher.Name = strconv.Itoa(i)
//         philosopher.DominantHand = fork0
//         philosopher.OtherHand = fork1
//         philosophers = append(philosophers, philosopher)
//         fork0 = fork1
//     }
//     philosopher := Philosopher{}
//     philosopher.Name = strconv.Itoa(0)
//     philosopher.DominantHand = forkLeft
//     philosopher.OtherHand = fork0
//     philosophers = append(philosophers, philosopher)

//     dt.Philosophers = philosophers
//     dt.isDone.Add(5)

// }

// func (dt *DiningTable) StartDining() {
//     for _,p := range dt.Philosophers {
//         rg := rGenerator(p.Name)
//         go p.Dine(rg, &dt.isDone)
//     }
// }

// func main()  {

//     var dt DiningTable

//     dt.Initialize()

//     dt.StartDining()

//     dt.isDone.Wait()
//     fmt.Println("dinner is done")
// }

// package main

// import "fmt"

// //Fibonacci series
// func Fibonacci(ch chan int, quit chan bool) {
//     x ,y := 0, 1
//     for {

//         select {
//         case ch <- x:
//         x, y = y, x+y
//         case <- quit:
//             fmt.Println("Fibonacci series concluded")
//             return
//         }
//     }

// }

// func main()  {
//     ch := make(chan int)
//     quit := make(chan bool)
//     var n int
//     fmt.Scanln(&n)
//     go func (n int)  {
//         for i := 0; i < n; i++ {
//             fmt.Println(<- ch)
//         }
//         quit <- true
//     }(n)

//     Fibonacci(ch, quit)
// }

// //generator and receiver functions
// func generator() chan int {
// 	ch := make(chan int)

// 	go func() {
// 		for i := 0; i < 10; i++ {
// 			ch <- i
// 		}
//         close(ch)
// 	}()

// 	return ch
// }

// func receiver(ch chan int) {
// 	for v := range ch {
// 		fmt.Println(v)
// 	}
// }

// func main() {
//     receiver(generator())
// }

// 10 goroutines
// func main() {
// 	ch := make(chan string)

// 	for i := 0; i < 10; i++ {
// 		go func(i int) {
// 			ch <- "go routine " + strconv.Itoa(i)
// 		}(i)
// 	}

//     for i:=0; i< 11; i++ {
//         fmt.Println(<-ch)
//     }
// }\
// 	start := time.Now();
// 	n := 1000*1000*10;
// 	var mx int64 =0;
// 	for i := 0; i<n; i++ {
// 		mx = max(mx, int64(n))
// 	}

//     elapsed := time.Since(start)
//     fmt.Println("Time taken:", elapsed)

// 	ch := make(chan int)
// 	var chh chan int
// 	chh = make(chan int)
// 	for i :=1; i< 10; i++ {
// 	go func(i int) {
// 		ch <- i
// 		chh <- i+5
// 		}(i)

// 	}

// 	fmt.Println(<-ch)
// 	fmt.Println(<-ch)
// 	fmt.Println(<-chh)

// for {
//     select {
//     case val, ok := <-ch:
//         if !ok {
//             fmt.Println("Channel closed, exiting.")
//             return
//         }
//         fmt.Println("Received:", val)
// 	case val, ok := <- chh:
// 		if !ok {
// 			fmt.Println("channel closed, existing.")
// 			return
// 		}
// 		fmt.Println("cchhhh ", val)
//     default:
//         fmt.Println("Waiting...")
//         time.Sleep(100 * time.Millisecond)
//     }
// }

// ticker := time.NewTicker(1 * time.Second)
// ticker.Reset(2* time.Second)
// wg := sync.WaitGroup{}
// n := 5;
// for {
// 	n--
// 	if(n <0) {
// 		break
// 	}
// 	select {
// 	case <-ticker.C:
// 		wg.Add(1)
// 		go printTime(&wg)
// 	}
// }
// wg.Wait()

func printTime(wg *sync.WaitGroup) int {
		fmt.Println(time.Now())
		wg.Done()
		return 2
}
func main()  {

  ch := make(chan string, 1)
  wg := &sync.WaitGroup{}
  
  a := "hello"
  wg.Add(1)
  go func (a string)  {
	// wg.Add(1)
	 ch <- a
	 fmt.Println("go routine 1 is done")
	//  wg.Done()
  }(a)

  go func() {
	// wg.Add(1)
	time.Sleep(2*time.Second)
   for {
	select {
	case val, ok := <-ch :
		if ok {
			fmt.Println(val + " world")
		}
		wg.Done()
		return

	default:
		fmt.Println("going to sleep")
		time.Sleep(1* time.Second)
	}
   }
  }()
   wg.Wait()
}

