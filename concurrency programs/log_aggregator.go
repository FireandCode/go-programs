package main

import (
	"fmt"
	"sync"
	"time"
)

/*
You have multiple log sources (producers), each emitting logs independently.
A central processor (consumer) collects and analyzes them for error patterns.
Goal:
Make log producers and consumers run concurrently and ensure no logs are missed or processed out of order beyond reason.

*/

func consumer(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for val := range ch {
		fmt.Println("consuming log ", val)
	}
	
}

func Producer(ch chan string, num int64, wgP *sync.WaitGroup) {
	for i := int64(0); i<10; i++ {
		ch <- fmt.Sprintf("Log %d from Producer %d", i,num)
		time.Sleep(10* time.Millisecond)
	}
	wgP.Done()
}

func runProducer(ch chan string) {
	wgP := &sync.WaitGroup{}
	
	for i :=int64(0); i<5; i++ {
		wgP.Add(1)
		go Producer(ch,i, wgP)
	}

	wgP.Wait()
	close(ch)
}

func main() {
	fmt.Println("hello log_aggregator")
	ch := make(chan string, 20)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go consumer(ch, wg)
	
	go runProducer(ch)

	wg.Wait()
}