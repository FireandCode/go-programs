/*
*Detailed Problem Statement
You are tasked with designing and implementing a concurrent job processing system in Go.
The system must be capable of handling a continuous stream of incoming jobs while efficiently
managing a pool of worker threads to process these jobs concurrently.
The system must satisfy the following requirements:
Job Queue:

	Incoming jobs must be placed into a queue as soon as they are received. The queue must have a finite capacity.
	If the queue is full, the system must appropriately handle the situation (e.g., by blocking, dropping, or rejecting new jobs).

Worker Pool:

	A fixed number of worker threads (goroutines) must be initialized at the start of the system.
	 These workers are responsible for pulling jobs from the queue and processing them. Each worker
	  should operate independently and continuously pull and process jobs until the system is shut down.

Concurrency Management:

	The system must efficiently coordinate the communication between the job queue and workers using concurrency
	 primitives available in Go.

Job Processing Simulation:

	Each job can be simulated as a task that takes a random amount of time (within a specified range) to complete.

System Shutdown:

	The system should provide a mechanism to gracefully shut down, ensuring that all queued jobs are processed
	 before termination and that all worker threads exit cleanly.

Metrics (Optional, Advanced):

	Optionally, the system can keep track of simple metrics, such as the number of completed jobs, the number of
	 failed jobs, or the average processing time per job.

Worker

	ID
	Execute()

Queue

	Push
	Pop()
	Get()

Jobs

	ID
	Task

JobExecutor

	workers []Worker
	queue Queue
	jobs []Job
*/
package main

import "fmt"



func main()  {
	a := 10;
	
	for i := 0; i < a; i++ {
		fmt.Println(i)
	}
}