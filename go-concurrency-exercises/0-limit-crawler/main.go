//////////////////////////////////////////////////////////////////////
//
// Your task is to change the code to limit the crawler to at most one
// page per second, while maintaining concurrency (in other words,
// Crawl() must be called concurrently)
//
// @hint: you can achieve this by adding 3 lines
//

package main

import (
	"fmt"
	"sync"
	"time"
)

// Crawl uses `fetcher` from the `mockfetcher.go` file to imitate a
// real crawler. It crawls until the maximum depth has reached.
func Crawl(url string, depth int, wg *sync.WaitGroup, rateLimiter <-chan time.Time) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	<- rateLimiter
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)

	wg.Add(len(urls))
	for _, u := range urls {
		// Do not remove the `go` keyword, as Crawl() must be
		// called concurrently
		go Crawl(u, depth-1, wg, rateLimiter)
	}
}

func main() {
	var wg sync.WaitGroup
	
	rateLimiter := time.NewTicker(1* time.Second)
	defer rateLimiter.Stop()
	wg.Add(1)
	Crawl("http://golang.org/", 4, &wg, rateLimiter.C)
	wg.Wait()
}
