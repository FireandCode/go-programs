//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	mu        sync.Mutex // For protecting channel initialization
	ch        chan struct{} // For serializing requests per user
}

func ProcessUsingChannel(process func(), ch chan struct{})  {
	process()
	close(ch)

}

// getOrCreateChannel ensures the user's channel is initialized
func (u *User) getOrCreateChannel() chan struct{} {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.ch == nil {
		u.ch = make(chan struct{}, 1) // Buffered channel with size 1
	}
	return u.ch
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {

	if u.IsPremium {
		process()
		return true
	}
	
	// Serialize requests for this user using channel
	userCh := u.getOrCreateChannel()
	userCh <- struct{}{} // Acquire the channel
	defer func() {
		<-userCh // Release the channel
	}()
	
	// Now only one request can run at a time for this user
	remainingTime := 10 - u.TimeUsed
	if remainingTime <= 0 {
		return false
	}

	start := time.Now()
	ch := make(chan struct{})
	go ProcessUsingChannel(process, ch)
	// go process()
	select {
	case <-time.After(time.Duration(remainingTime) * time.Second):
		u.TimeUsed += remainingTime
		return false
	case <-ch:
		elapsedTime := time.Since(start).Seconds()
		u.TimeUsed += int64(elapsedTime)
		return true
	}
}

func main() {
	RunMockServer()
}
