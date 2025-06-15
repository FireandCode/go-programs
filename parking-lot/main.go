package main

import (
	"fmt"
	"time"
)

/*

parking lot - struct.
	- allocateASpot
	- freeUpASpot
	- spotAvailable
Vehicles - interface
	- truck - struct
	- motorcycle - struct
	- car - struct

parkingSpot - struct
	id
	type ()
		- compact
		- large
		- small
		- electric
	isOccupied

ticket - struct
	- vehicle no.
	- person name
	- bill amount
	- spot id
	- entrytime
	- exittime


EntryGate
	- generateTicket //call the parkingLot.spotAvailable

ExitGate
	- processPayment //call the parkingLot.freeUpSpot

*/

// 🚗 VEHICLE INTERFACE
type Vehicle interface {
	getModel() string
	setModel(string)
	getNumberPlate() string
	setNumberPlate(string)
}

// 🚘 BASE VEHICLE STRUCT
type BaseVehicle struct {
	model       string
	numberPlate string
}

func (v *BaseVehicle) getModel() string             { return v.model }
func (v *BaseVehicle) setModel(model string)        { v.model = model }
func (v *BaseVehicle) getNumberPlate() string       { return v.numberPlate }
func (v *BaseVehicle) setNumberPlate(numberPlate string) { v.numberPlate = numberPlate }

// 🚗 SPECIFIC VEHICLES
type Car struct{ BaseVehicle }
type Truck struct{ BaseVehicle }
type MotorCycle struct{ BaseVehicle }

// 🅿️ PARKING SPOT TYPE ENUM
type SpotType int

const (
	Compact SpotType = iota
	Large
	Small
	Electric
)

// 🅿️ PARKING SPOT STRUCT
type ParkingSpot struct {
	id        int
	spotType  SpotType
	isOccupied bool
}

// 🎫 PARKING TICKET STRUCT
type Ticket struct {
	vehicleNumber string
	personName    string
	billAmount    float64
	spotID        int
	entryTime     time.Time
	exitTime      time.Time
}

// 🅿️ PARKING LOT STRUCT (Stores Tickets)
type ParkingLot struct {
	spots   []ParkingSpot
	tickets map[string]*Ticket // Store tickets using vehicle number as key
}

// Check if a spot is available
func (p *ParkingLot) spotAvailable(spotType SpotType) (*ParkingSpot, bool) {
	for i, spot := range p.spots {
		if spot.spotType == spotType && !spot.isOccupied {
			return &p.spots[i], true
		}
	}
	return nil, false
}

// Allocate a parking spot
func (p *ParkingLot) allocateASpot(spotType SpotType) (*ParkingSpot, bool) {
	spot, available := p.spotAvailable(spotType)
	if available {
		spot.isOccupied = true
		return spot, true
	}
	return nil, false
}

// Free up a parking spot
func (p *ParkingLot) freeUpASpot(spotID int) bool {
	for i := range p.spots {
		if p.spots[i].id == spotID {
			p.spots[i].isOccupied = false
			return true
		}
	}
	return false
}

// Store ticket
func (p *ParkingLot) addTicket(ticket *Ticket) {
	p.tickets[ticket.vehicleNumber] = ticket
}

// Retrieve ticket
func (p *ParkingLot) getTicket(vehicleNumber string) (*Ticket, bool) {
	ticket, exists := p.tickets[vehicleNumber]
	return ticket, exists
}

// Delete ticket
func (p *ParkingLot) removeTicket(vehicleNumber string) {
	delete(p.tickets, vehicleNumber)
}

// 🚪 ENTRY GATE STRUCT
type EntryGate struct {
	parkingLot *ParkingLot
}

// Generate a ticket when a vehicle enters
func (e *EntryGate) generateTicket(vehicle Vehicle, spotType SpotType) (*Ticket, bool) {
	spot, allocated := e.parkingLot.allocateASpot(spotType)
	if !allocated {
		fmt.Println("No available spots!")
		return nil, false
	}

	newTicket := &Ticket{
		vehicleNumber: vehicle.getNumberPlate(),
		personName:    "John Doe", // Placeholder
		billAmount:    0,
		spotID:        spot.id,
		entryTime:     time.Now(),
	}

	// Store ticket
	e.parkingLot.addTicket(newTicket)

	fmt.Printf("Ticket Generated: Vehicle %s parked at Spot %d\n", vehicle.getNumberPlate(), spot.id)
	return newTicket, true
}

// 🚪 EXIT GATE STRUCT
type ExitGate struct {
	parkingLot *ParkingLot
}

// Process payment when a vehicle exits
func (e *ExitGate) processPayment(vehicleNumber string) bool {
	ticket, exists := e.parkingLot.getTicket(vehicleNumber)
	if !exists {
		fmt.Println("Invalid ticket")
		return false
	}

	// Billing logic (₹10 per hour)
	duration := time.Since(ticket.entryTime).Hours()
	ticket.billAmount = duration * 10 // ₹10 per hour

	// Free up the parking spot
	if e.parkingLot.freeUpASpot(ticket.spotID) {
		// Remove ticket from storage
		e.parkingLot.removeTicket(vehicleNumber)
		fmt.Printf("Payment of ₹%.2f processed for Vehicle %s\n", ticket.billAmount, ticket.vehicleNumber)
		return true
	}
	return false
}

// 🎯 MAIN FUNCTION: Simulating the Parking Lot System
func main() {
	// Initialize Parking Lot with 5 spots and a ticket storage
	parkingLot := ParkingLot{
		spots: []ParkingSpot{
			{1, Compact, false},
			{2, Large, false},
			{3, Small, false},
			{4, Electric, false},
			{5, Compact, false},
		},
		tickets: make(map[string]*Ticket),
	}

	entryGate := EntryGate{&parkingLot}
	exitGate := ExitGate{&parkingLot}

	// Create a new vehicle
	myCar := Car{}
	myCar.setModel("Tesla Model 3")
	myCar.setNumberPlate("ABC123")

	// Vehicle enters parking lot
	_, success := entryGate.generateTicket(&myCar, Compact)
	if success {
		time.Sleep(2 * time.Second) // Simulating time spent in parking

		// Vehicle exits parking lot
		exitGate.processPayment(myCar.getNumberPlate())
	}
}
