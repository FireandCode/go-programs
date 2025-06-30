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

functional requirements
1. create a queue which will store tasks.
2. initalize workers goroutines at the start of the system. this workers will pick the jobs from queue and execute them.
3. graceful shutdown
4. add metrics as well

non functional requirements
1. concurrency

how should we handle the jobs which are coming with queue already filled.
-> 1. having a backup queue
-> 2. 3 retries with 2 second gap betweent hem
2 approach is fine.

	type Task {
		Id string
		message string
	}

	type struct Node {
		task Task
		next Node*
	}

	type struct queue {
		capacity int
		head *Node
		tail *Node

		push(Task) (string, err)
		pop() (Task, err)
	}
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)
type Task struct {
	Id string
	message string
}

type Worker struct {
	ID string
	ctx context.Context
	wg *sync.WaitGroup
	metric *Metric
}

type Metric struct {
	noOfJobs    int64
	successJobs int64
	failedJobs  int64
}

type ThreadPool struct{
	jobs chan Task // Buffered channel for jobs
	workers []Worker
	ctx context.Context
	metric Metric
	wg sync.WaitGroup
}


func(w *Worker) process(jobs <-chan Task)  {
	for task := range jobs {
		w.doTask(task)
	}
	fmt.Printf("worker %s- is stopped\n", w.ID)
}

func (w *Worker) doTask(task Task) {
	wait := rand.Intn(3) + 1
	time.Sleep(time.Duration(wait) * time.Second)
	// Simulate random failure (20% chance)
	if rand.Float64() < 0.2 {
		atomic.AddInt64(&w.metric.failedJobs, 1)
	} else {
		atomic.AddInt64(&w.metric.successJobs, 1)
	}
	w.wg.Done()
}

func generateID()  string {
	return strconv.Itoa(rand.Int())
}

func InitializeThreadPool(capacity int, noOfWorkers int, ctx context.Context) *ThreadPool {
	threadPool := &ThreadPool{
		jobs: make(chan Task, capacity), // Buffered channel
		workers: make([]Worker, 0),
		ctx: ctx,
		metric: Metric{
			noOfJobs: 0,
			successJobs: 0,
			failedJobs: 0,
		},
	}

	for i := 0; i < noOfWorkers; i++ {
		id := generateID()
		worker := Worker{
			ID: id,
			ctx: ctx,
			wg: &threadPool.wg,
			metric: &threadPool.metric,
		}
		go worker.process(threadPool.jobs)
		threadPool.workers = append(threadPool.workers, worker)
	}

	return threadPool
}

func(th *ThreadPool) pushTask(message string) (Task, error) {
	task := Task{
		Id: generateID(),
		message: message,
	}
	for i:=0; i<3; i++ {
		select {
		case th.jobs <- task:
			th.wg.Add(1)
			atomic.AddInt64(&th.metric.noOfJobs, 1)
			return task, nil
		default:
			time.Sleep(2*time.Second)
		}
	}
	return Task{}, errors.New("Task is not pushed")
}



func(th *ThreadPool) Stop(cancel context.CancelFunc)  {
	close(th.jobs) // Signal workers to stop after all jobs are processed
	cancel() // Cancel context if needed
}


func main()  {
	rand.Seed(time.Now().UnixNano())
	ctx, cancel := context.WithCancel(context.Background())
	th := InitializeThreadPool(10, 3, ctx)
	fmt.Println("HELLLO")
	submitWg := sync.WaitGroup{}
	for i:=0; i<30; i++ {
		submitWg.Add(1)
		go func(i int) {
			defer submitWg.Done()
			fmt.Printf("task %d", i)
			res, err := th.pushTask(fmt.Sprintf("task %d", i))
			if err != nil {
				fmt.Printf("task %d is not successful pushed to queue", i)
				return
			}
			fmt.Println("success: task ", res.Id)
		}(i)
	}
	submitWg.Wait() // Wait for all jobs to be submitted
	th.wg.Wait()    // Wait for all jobs to be processed
	th.Stop(cancel)
	// Print metrics summary
	fmt.Println("\n--- Metrics ---")
	fmt.Printf("Total jobs submitted: %d\n", atomic.LoadInt64(&th.metric.noOfJobs))
	fmt.Printf("Successful jobs: %d\n", atomic.LoadInt64(&th.metric.successJobs))
	fmt.Printf("Failed jobs: %d\n", atomic.LoadInt64(&th.metric.failedJobs))
}