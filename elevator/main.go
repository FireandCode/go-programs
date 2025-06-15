package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

// Direction Enum
type Direction int

const (
	UP Direction = iota
	DOWN
	IDLE_DIR
)

// Status Enum
type Status int

const (
	IDLE Status = iota
	MOVING
	STOPPED
)

// Passenger Struct
type Passenger struct {
	name       string
	entryFloor int
	exitFloor  int
}

// FloorPriorityQueue (Min-Heap for LOOK Algorithm)
type FloorPriorityQueue []int

func (pq FloorPriorityQueue) Len() int           { return len(pq) }
func (pq FloorPriorityQueue) Less(i, j int) bool { return pq[i] < pq[j] } // Min-Heap
func (pq FloorPriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *FloorPriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(int))
}
func (pq *FloorPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// Elevator Struct
type Elevator struct {
	id            int
	currentFloor  int
	status        Status
	direction     Direction
	passengers    []Passenger
	floorsToStop  *FloorPriorityQueue
	mutex         sync.Mutex
}

// NewElevator Constructor
func NewElevator(id int) *Elevator {
	pq := &FloorPriorityQueue{}
	heap.Init(pq)
	return &Elevator{id: id, currentFloor: 0, status: IDLE, direction: IDLE_DIR, floorsToStop: pq}
}

// Add Floor Request (LOOK Algorithm)
func (e *Elevator) AddFloorRequest(floor int) {
	e.mutex.Lock()
	heap.Push(e.floorsToStop, floor)
	e.mutex.Unlock()
}

// Elevator Movement (LOOK Algorithm)
func (e *Elevator) MoveElevator() {
	for {
		e.mutex.Lock()
		if e.floorsToStop.Len() == 0 {
			e.status = IDLE
			e.direction = IDLE_DIR
			e.mutex.Unlock()
			time.Sleep(2 * time.Second) // Wait before checking again
			continue
		}

		// Determine direction
		nextFloor := (*e.floorsToStop)[0]
		if nextFloor > e.currentFloor {
			e.direction = UP
		} else if nextFloor < e.currentFloor {
			e.direction = DOWN
		}

		// Move towards next floor
		e.status = MOVING
		fmt.Printf("Elevator %d moving %v\n", e.id, e.direction)
		for e.currentFloor != nextFloor {
			time.Sleep(1 * time.Second) // Simulate movement time
			if e.direction == UP {
				e.currentFloor++
			} else {
				e.currentFloor--
			}
			fmt.Printf("Elevator %d at floor %d\n", e.id, e.currentFloor)
		}

		// Stop at the floor
		e.status = STOPPED
		fmt.Printf("Elevator %d stopped at floor %d\n", e.id, e.currentFloor)

		// Remove the floor from queue
		heap.Pop(e.floorsToStop)
		e.mutex.Unlock()
		time.Sleep(2 * time.Second) // Simulate stop time
	}
}

// Building Struct
type Building struct {
	elevators []*Elevator
}

// NewBuilding Constructor
func NewBuilding(numElevators int) *Building {
	elevators := make([]*Elevator, numElevators)
	for i := 0; i < numElevators; i++ {
		elevators[i] = NewElevator(i + 1)
	}
	return &Building{elevators: elevators}
}

// Assign Request to the Best Elevator
func (b *Building) RequestElevator(floor int) {
	// Select the best elevator (closest idle/moving in the same direction)
	bestElevator := b.elevators[0]
	for _, elevator := range b.elevators {
		if elevator.status == IDLE || (elevator.direction == UP && floor > elevator.currentFloor) || (elevator.direction == DOWN && floor < elevator.currentFloor) {
			bestElevator = elevator
			break
		}
	}

	fmt.Printf("Request for floor %d assigned to Elevator %d\n", floor, bestElevator.id)
	bestElevator.AddFloorRequest(floor)
}

func main() {
	// Create a Building with 2 Elevators
	building := NewBuilding(2)

	// Start elevator movement goroutines
	for _, elevator := range building.elevators {
		go elevator.MoveElevator()
	}

	// Simulate Requests
	time.Sleep(1 * time.Second)
	building.RequestElevator(3)
	time.Sleep(2 * time.Second)
	building.RequestElevator(7)
	time.Sleep(3 * time.Second)
	building.RequestElevator(2)

	// Keep the main routine alive
	select {}
}
