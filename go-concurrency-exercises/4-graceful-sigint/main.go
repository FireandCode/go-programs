//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

func main() {
	// Create a process
	proc := MockProcess{}
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, syscall.SIGINT)
	count := int32(0)
	
	// Run the process in a goroutine
	go proc.Run()
	
	// Handle signals in main goroutine
	for sig := range ch {
		fmt.Println("Received Signal and count:", sig, count)
		if atomic.AddInt32(&count, 1) == 1 {
			fmt.Println("\nReceived Signal:", sig)
			fmt.Println("Stopping the process gracefully...")
			go proc.Stop() // Run Stop() in a goroutine so we can handle more signals
		} else {
			fmt.Println("\nReceived second Signal:", sig)
			fmt.Println("Force killing the program...")
			os.Exit(1)
		}
	}
}
