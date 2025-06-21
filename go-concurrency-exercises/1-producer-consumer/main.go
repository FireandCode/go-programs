//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, chTw chan *Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(chTw)
			return
		}
			chTw <- tweet
	}
}

func consumer(chTw chan *Tweet, wg *sync.WaitGroup) {
	for t := range chTw {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}

	wg.Done()
}

func main() {
	start := time.Now()
	stream := GetMockStream()
	chTw := make(chan *Tweet, 6)
	wg := &sync.WaitGroup{}
	// Producer
	
	go producer(stream, chTw)

	// Consumer
	for i:=0; i<3; i++ {
	wg.Add(1)
	go consumer(chTw, wg)

	}
	
	wg.Wait()
	fmt.Printf("Process took %s\n", time.Since(start))
}
