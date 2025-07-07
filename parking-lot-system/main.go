package main

import (
	"fmt"
	"sync"
	"time"
)

// --- ENUMS & CONSTANTS ---
type VehicleType string
type SlotStatus string
type TicketStatus string

const (
	VehicleTypeCar   VehicleType = "CAR"
	VehicleTypeBike  VehicleType = "BIKE"
	VehicleTypeTruck VehicleType = "TRUCK"

	SlotStatusAvailable SlotStatus = "AVAILABLE"
	SlotStatusOccupied  SlotStatus = "OCCUPIED"
	SlotStatusReserved  SlotStatus = "RESERVED"

	TicketStatusActive   TicketStatus = "ACTIVE"
	TicketStatusCompleted TicketStatus = "COMPLETED"
)

// --- VEHICLE (Factory Pattern) ---
type Vehicle struct {
	ID         string
	Type       VehicleType
	LicensePlate string
	EntryTime  time.Time
	ExitTime   time.Time
}

func NewVehicle(vehicleType VehicleType, licensePlate string) *Vehicle {
	return &Vehicle{
		ID:           generateID(),
		Type:         vehicleType,
		LicensePlate: licensePlate,
		EntryTime:    time.Now(),
	}
}

// --- SLOT ---
type Slot struct {
	ID       string
	FloorID  string
	Number   int
	Type     VehicleType
	Status   SlotStatus
	Vehicle  *Vehicle
	mu       sync.Mutex
}

func NewSlot(floorID string, number int, slotType VehicleType) *Slot {
	return &Slot{
		ID:      generateID(),
		FloorID: floorID,
		Number:  number,
		Type:    slotType,
		Status:  SlotStatusAvailable,
	}
}

func (s *Slot) ParkVehicle(vehicle *Vehicle) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.Status == SlotStatusAvailable && s.Type == vehicle.Type {
		s.Vehicle = vehicle
		s.Status = SlotStatusOccupied
		return true
	}
	return false
}

func (s *Slot) RemoveVehicle() *Vehicle {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.Status == SlotStatusOccupied {
		vehicle := s.Vehicle
		s.Vehicle = nil
		s.Status = SlotStatusAvailable
		return vehicle
	}
	return nil
}

// --- TICKET ---
type Ticket struct {
	ID        string
	Vehicle   *Vehicle
	Slot      *Slot
	EntryTime time.Time
	ExitTime  time.Time
	Status    TicketStatus
	Amount    float64
}

func NewTicket(vehicle *Vehicle, slot *Slot) *Ticket {
	return &Ticket{
		ID:        generateID(),
		Vehicle:   vehicle,
		Slot:      slot,
		EntryTime: time.Now(),
		Status:    TicketStatusActive,
	}
}

func (t *Ticket) CalculateAmount() float64 {
	if t.ExitTime.IsZero() {
		t.ExitTime = time.Now()
	}
	duration := t.ExitTime.Sub(t.EntryTime)
	hours := duration.Hours()
	
	// Simple pricing: $2 per hour for cars, $1 for bikes, $5 for trucks
	rate := 2.0
	switch t.Vehicle.Type {
	case VehicleTypeBike:
		rate = 1.0
	case VehicleTypeTruck:
		rate = 5.0
	}
	
	t.Amount = hours * rate
	return t.Amount
}

// --- FLOOR ---
type Floor struct {
	ID    string
	Name  string
	Slots []*Slot
	mu    sync.RWMutex
}

func NewFloor(id, name string) *Floor {
	return &Floor{
		ID:    id,
		Name:  name,
		Slots: make([]*Slot, 0),
	}
}

func (f *Floor) AddSlot(slot *Slot) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Slots = append(f.Slots, slot)
}

func (f *Floor) GetAvailableSlots(vehicleType VehicleType) []*Slot {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	var available []*Slot
	for _, slot := range f.Slots {
		if slot.Status == SlotStatusAvailable && slot.Type == vehicleType {
			available = append(available, slot)
		}
	}
	return available
}

// --- ALLOCATION STRATEGY (Strategy Pattern) ---
type AllocationStrategy interface {
	AllocateSlot(floors []*Floor, vehicleType VehicleType) (*Slot, *Floor)
}

// Nearest Slot Strategy
type NearestSlotStrategy struct{}

func (s *NearestSlotStrategy) AllocateSlot(floors []*Floor, vehicleType VehicleType) (*Slot, *Floor) {
	for _, floor := range floors {
		availableSlots := floor.GetAvailableSlots(vehicleType)
		if len(availableSlots) > 0 {
			return availableSlots[0], floor // Return first available slot
		}
	}
	return nil, nil
}

// Type-Based Strategy (prefers specific floors for specific vehicle types)
type TypeBasedStrategy struct{}

func (s *TypeBasedStrategy) AllocateSlot(floors []*Floor, vehicleType VehicleType) (*Slot, *Floor) {
	// For trucks, prefer ground floor (index 0)
	// For bikes, prefer higher floors
	// For cars, any floor
	
	preferredFloors := make([]*Floor, len(floors))
	copy(preferredFloors, floors)
	
	switch vehicleType {
	case VehicleTypeTruck:
		// Trucks prefer ground floor
		if len(preferredFloors) > 0 {
			availableSlots := preferredFloors[0].GetAvailableSlots(vehicleType)
			if len(availableSlots) > 0 {
				return availableSlots[0], preferredFloors[0]
			}
		}
	case VehicleTypeBike:
		// Bikes prefer higher floors (reverse order)
		for i := len(preferredFloors) - 1; i >= 0; i-- {
			availableSlots := preferredFloors[i].GetAvailableSlots(vehicleType)
			if len(availableSlots) > 0 {
				return availableSlots[0], preferredFloors[i]
			}
		}
	}
	
	// Fallback to nearest strategy
	nearestStrategy := &NearestSlotStrategy{}
	return nearestStrategy.AllocateSlot(floors, vehicleType)
}

// --- GATE ---
type Gate struct {
	ID   string
	Name string
	mu   sync.Mutex
}

func NewGate(id, name string) *Gate {
	return &Gate{
		ID:   id,
		Name: name,
	}
}

// --- PARKING LOT (Main Entity) ---
type ParkingLot struct {
	ID           string
	Name         string
	Floors       []*Floor
	Gates        []*Gate
	Strategy     AllocationStrategy
	Tickets      map[string]*Ticket
	mu           sync.RWMutex
}

func NewParkingLot(id, name string) *ParkingLot {
	return &ParkingLot{
		ID:       id,
		Name:     name,
		Floors:   make([]*Floor, 0),
		Gates:    make([]*Gate, 0),
		Strategy: &NearestSlotStrategy{}, // Default strategy
		Tickets:  make(map[string]*Ticket),
	}
}

func (pl *ParkingLot) AddFloor(floor *Floor) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	pl.Floors = append(pl.Floors, floor)
}

func (pl *ParkingLot) AddGate(gate *Gate) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	pl.Gates = append(pl.Gates, gate)
}

func (pl *ParkingLot) SetStrategy(strategy AllocationStrategy) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	pl.Strategy = strategy
}

// Park Vehicle - Main parking logic
func (pl *ParkingLot) ParkVehicle(vehicle *Vehicle) (*Ticket, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	
	// Find available slot using strategy
	slot, floor := pl.Strategy.AllocateSlot(pl.Floors, vehicle.Type)
	if slot == nil {
		return nil, fmt.Errorf("no available slot for vehicle type %s", vehicle.Type)
	}
	
	// Park the vehicle
	if !slot.ParkVehicle(vehicle) {
		return nil, fmt.Errorf("failed to park vehicle in slot %s", slot.ID)
	}
	
	// Create ticket
	ticket := NewTicket(vehicle, slot)
	pl.Tickets[ticket.ID] = ticket
	
	fmt.Printf("Vehicle %s parked in Floor %s, Slot %d\n", vehicle.LicensePlate, floor.Name, slot.Number)
	return ticket, nil
}

// Exit Vehicle
func (pl *ParkingLot) ExitVehicle(ticketID string) (*Ticket, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	
	ticket, exists := pl.Tickets[ticketID]
	if !exists {
		return nil, fmt.Errorf("ticket %s not found", ticketID)
	}
	
	if ticket.Status == TicketStatusCompleted {
		return nil, fmt.Errorf("ticket %s already completed", ticketID)
	}
	
	// Remove vehicle from slot
	vehicle := ticket.Slot.RemoveVehicle()
	if vehicle == nil {
		return nil, fmt.Errorf("failed to remove vehicle from slot")
	}
	
	// Calculate amount
	amount := ticket.CalculateAmount()
	ticket.Status = TicketStatusCompleted
	vehicle.ExitTime = time.Now()
	
	fmt.Printf("Vehicle %s exited. Amount: $%.2f\n", vehicle.LicensePlate, amount)
	return ticket, nil
}

// Get parking status
func (pl *ParkingLot) GetStatus() map[string]interface{} {
	pl.mu.RLock()
	defer pl.mu.RUnlock()
	
	status := make(map[string]interface{})
	status["total_floors"] = len(pl.Floors)
	status["total_gates"] = len(pl.Gates)
	status["active_tickets"] = len(pl.Tickets)
	
	// Count available slots by type
	availableSlots := make(map[VehicleType]int)
	totalSlots := make(map[VehicleType]int)
	
	for _, floor := range pl.Floors {
		for _, slot := range floor.Slots {
			totalSlots[slot.Type]++
			if slot.Status == SlotStatusAvailable {
				availableSlots[slot.Type]++
			}
		}
	}
	
	status["available_slots"] = availableSlots
	status["total_slots"] = totalSlots
	
	return status
}

// --- UTILS ---
var idCounter = 0
var idMu sync.Mutex

func generateID() string {
	idMu.Lock()
	defer idMu.Unlock()
	idCounter++
	return fmt.Sprintf("id-%d", idCounter)
}

// --- MAIN: Demonstrate the system ---
func main() {
	// Create parking lot
	parkingLot := NewParkingLot("PL001", "Downtown Parking")
	
	// Add floors
	groundFloor := NewFloor("F1", "Ground Floor")
	firstFloor := NewFloor("F2", "First Floor")
	secondFloor := NewFloor("F3", "Second Floor")
	
	// Add slots to ground floor (mostly trucks and cars)
	for i := 1; i <= 10; i++ {
		groundFloor.AddSlot(NewSlot("F1", i, VehicleTypeCar))
	}
	for i := 11; i <= 15; i++ {
		groundFloor.AddSlot(NewSlot("F1", i, VehicleTypeTruck))
	}
	
	// Add slots to first floor (cars and bikes)
	for i := 1; i <= 20; i++ {
		firstFloor.AddSlot(NewSlot("F2", i, VehicleTypeCar))
	}
	for i := 21; i <= 30; i++ {
		firstFloor.AddSlot(NewSlot("F2", i, VehicleTypeBike))
	}
	
	// Add slots to second floor (mostly bikes)
	for i := 1; i <= 25; i++ {
		secondFloor.AddSlot(NewSlot("F3", i, VehicleTypeBike))
	}
	
	parkingLot.AddFloor(groundFloor)
	parkingLot.AddFloor(firstFloor)
	parkingLot.AddFloor(secondFloor)
	
	// Add gates
	entryGate := NewGate("G1", "Entry Gate")
	exitGate := NewGate("G2", "Exit Gate")
	parkingLot.AddGate(entryGate)
	parkingLot.AddGate(exitGate)
	
	// Test with different strategies
	fmt.Println("=== Testing Nearest Slot Strategy ===")
	parkingLot.SetStrategy(&NearestSlotStrategy{})
	
	// Park some vehicles
	car1 := NewVehicle(VehicleTypeCar, "ABC123")
	ticket1, _ := parkingLot.ParkVehicle(car1)
	
	bike1 := NewVehicle(VehicleTypeBike, "XYZ789")
	ticket2, _ := parkingLot.ParkVehicle(bike1)
	
	truck1 := NewVehicle(VehicleTypeTruck, "TRK456")
	ticket3, _ := parkingLot.ParkVehicle(truck1)
	
	// Switch to type-based strategy
	fmt.Println("\n=== Testing Type-Based Strategy ===")
	parkingLot.SetStrategy(&TypeBasedStrategy{})
	
	car2 := NewVehicle(VehicleTypeCar, "DEF456")
	ticket4, _ := parkingLot.ParkVehicle(car2)
	
	bike2 := NewVehicle(VehicleTypeBike, "MNO123")
	ticket5, _ := parkingLot.ParkVehicle(bike2)
	
	// Exit some vehicles
	parkingLot.ExitVehicle(ticket1.ID)
	parkingLot.ExitVehicle(ticket2.ID)
	
	// Show status
	fmt.Println("\n=== Parking Lot Status ===")
	status := parkingLot.GetStatus()
	fmt.Printf("Total Floors: %d\n", status["total_floors"])
	fmt.Printf("Total Gates: %d\n", status["total_gates"])
	fmt.Printf("Active Tickets: %d\n", status["active_tickets"])
	fmt.Printf("Available Slots: %v\n", status["available_slots"])
	fmt.Printf("Total Slots: %v\n", status["total_slots"])
} 