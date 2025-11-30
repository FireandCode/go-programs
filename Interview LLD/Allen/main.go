package main

import (
	"errors"
	"fmt"
)

/*










 Additional Requirements
 1. vehicle priority - done
 2. slot reallocation - will go in ClearASlot
 3. Fee Calculation - done

 Scalability
- support concurrency
*/

type ParkingLot struct {
	slots *[]Slot
	tickets []Ticket
	vehicles []Vehicle
}


type VehicleType int
	

const (
	BUS VehicleType = iota
	CAR
	MOTORCYCLE
)

/*
Slot struct
type (Enum - BUS, CAR, MOTORCYCLE)
isOccupied bool
registrationNumber string
isOccupiedByHigherPriority bool
*/
type Slot struct {
	id 						   int
	vehicleType                VehicleType
	isOccupied                 bool
	registrationNumber         string
	isOccupiedByHigherPriority bool
}

/*

Ticket struct
	- registration number
	- slot ID
	- amount
	- entryTime
	- exitTime
*/
type Ticket struct {
	 registrationNumber string
	 slotId int
	 amount int
	 entryTime int
	 exitTime int 
}

/*
Vehicle interface
ParkVehicle
RemoveVehicle
 BaseVehicle struct
- Bus - struct
- Car - struct
- Motocycle - struct
*/


type Vehicle interface {
	ParkVehicle(slots *[]Slot) (Slot, error)
	RemoveVehicle(slots *[]Slot, slot Slot) (bool, error)
	getRegistrationNumber() string
}

type BaseVehicle struct {
	vehicleType        VehicleType
	registrationNumber string
}

func (b *BaseVehicle) getVehicleType() VehicleType {
	return b.vehicleType
}

func (b *BaseVehicle) getRegistrationNumber() string {
	return b.registrationNumber
}

func (b *BaseVehicle) setVehicleType(vehicleType VehicleType) {
	b.vehicleType = vehicleType
}

func (b *BaseVehicle) setRegistrationNUmber(registrationNumber string) {
	b.registrationNumber = registrationNumber
}

func (b *BaseVehicle) ParkVehicle(slots *[]Slot) (Slot, error) {
	return Slot{}, nil
}

func (b *BaseVehicle) RemoveVehicle(slots *[]Slot, slot Slot) (bool, error) {
	return true, nil
}

type Motorcycle struct {
	BaseVehicle
}

func (m *Motorcycle) ParkVehicle(slots *[]Slot) (Slot, error) {

	for i := 0; i < len(*slots); i++ {
		if (*slots)[i].isOccupied == false &&
			(*slots)[i].vehicleType == m.vehicleType {
			slot := Slot{
				id: 					i+1,
				vehicleType:                m.vehicleType,
				registrationNumber:         m.registrationNumber,
				isOccupied:                 true,
				isOccupiedByHigherPriority: false,
			}
			(*slots)[i] = slot

			return slot, nil
		}
	}

	err := errors.New("No Available Slots")
	return Slot{}, err
}

func (m *Motorcycle) RemoveVehicle(slots *[]Slot, slot Slot) (bool, error) {

	for i := 0; i < len(*slots); i++ {
		if (*slots)[i].isOccupied == true &&
			(*slots)[i].id == slot.id {
			slot := Slot{
				id: 					i+1,
				vehicleType:                m.vehicleType,
				registrationNumber:         "",
				isOccupied:                 false,
				isOccupiedByHigherPriority: false,
			}
			(*slots)[i] = slot

			return true, nil
		}
	}

	err := errors.New("Slot id is not correct")
	return false, err
}

type Car struct {
	BaseVehicle
}

func (m *Car) ParkVehicle(slots *[]Slot) (Slot, error) {

	for i := 0; i < len(*slots); i++ {
		if (*slots)[i].isOccupied == false &&
			(*slots)[i].vehicleType == m.vehicleType {
			slot := Slot{
				id: 					i+1,
				vehicleType:                m.vehicleType,
				registrationNumber:         m.registrationNumber,
				isOccupied:                 true,
				isOccupiedByHigherPriority: false,
			}
			(*slots)[i] = slot

			return slot, nil
		}
	}

	//check the motorCycle Slots

	for i := 0; i < len(*slots); i++ {
		if (*slots)[i].isOccupied == false &&
			(*slots)[i].vehicleType == MOTORCYCLE {
			slot := Slot{
				id: 					i+1,
				vehicleType:                MOTORCYCLE,
				registrationNumber:         m.registrationNumber,
				isOccupied:                 true,
				isOccupiedByHigherPriority: true,
			}
			(*slots)[i] = slot
			fmt.Printf("Booked a motorcycle slot %d for car %s", i, m.registrationNumber )
			return slot, nil
		}
	}

	err := errors.New("No Available Slots")
	return Slot{}, err
}

func (m *Car) RemoveVehicle(slots *[]Slot, slot Slot) (bool, error) {

	for i := 0; i < len(*slots); i++ {
		if (*slots)[i].isOccupied == true &&
			(*slots)[i].id == slot.id {
			slot := Slot{
				id: 					i+1,
				vehicleType:                (*slots)[i].vehicleType,
				registrationNumber:         "",
				isOccupied:                 false,
				isOccupiedByHigherPriority: false,
			}
			(*slots)[i] = slot

			return true, nil
		}
	}

	err := errors.New("Slot id is not correct")
	return false, err
}

type Bus struct {
	BaseVehicle
}

func (m *Bus) ParkVehicle(slots *[]Slot) (Slot, error) {

	for i := 0; i < len(*slots); i++ {
		if (*slots)[i].isOccupied == false &&
			(*slots)[i].vehicleType == m.vehicleType {
			slot := Slot{
				id: 					i+1,
				vehicleType:                m.vehicleType,
				registrationNumber:         m.registrationNumber,
				isOccupied:                 true,
				isOccupiedByHigherPriority: false,
			}
			(*slots)[i] = slot

			return slot, nil
		}
	}

	//check the motorCycle Slots
	n := len(*slots)
	for i := 0; i < len(*slots); i++ {
		if (*slots)[i].isOccupied == false &&
			(*slots)[i].vehicleType == CAR &&
			i+1 <n && (*slots)[i+2].isOccupied == false &&
			(*slots)[i].vehicleType == CAR &&
			i+2 <n && (*slots)[i+2].isOccupied == false &&
			(*slots)[i].vehicleType == CAR {
			slot := Slot{
				id: 					i+1,
				vehicleType:                CAR,
				registrationNumber:         m.registrationNumber,
				isOccupied:                 true,
				isOccupiedByHigherPriority: true,
			}

			slot2 := Slot{
				id: 					i+2,
				vehicleType:                CAR,
				registrationNumber:         m.registrationNumber,
				isOccupied:                 true,
				isOccupiedByHigherPriority: true,
			}

			slot3 := Slot{
				id: 					i+3,
				vehicleType:                CAR,
				registrationNumber:         m.registrationNumber,
				isOccupied:                 true,
				isOccupiedByHigherPriority: true,
			}
			(*slots)[i] = slot
			(*slots)[i+1] = slot2
			(*slots)[i+2] = slot3
			fmt.Printf("Booked a bus slot %d for bus %s", i, m.registrationNumber )
			return slot, nil
		}
	}

	err := errors.New("No Available Slots")
	return Slot{}, err
}

func (m *Bus) RemoveVehicle(slots *[]Slot, slot Slot) (bool, error) {

	for i := 0; i < len(*slots); i++ {
		if (*slots)[i].isOccupied == true &&
			(*slots)[i].id == slot.id {
				if (*slots)[i].vehicleType == BUS  {
					slot := Slot{
						id: 					i+1,
						vehicleType:                (*slots)[i].vehicleType,
						registrationNumber:         "",
						isOccupied:                 false,
						isOccupiedByHigherPriority: false,
					}
					(*slots)[i] = slot
		
					return true, nil
				}

				slot := Slot{
					id: 					i+1,
					vehicleType:                CAR,
					registrationNumber:         "",
					isOccupied:                 false,
					isOccupiedByHigherPriority: false,
				}
	
				slot2 := Slot{
					id: 					i+2,
					vehicleType:                CAR,
					registrationNumber:         "",
					isOccupied:                 false,
					isOccupiedByHigherPriority: false,
				}
	
				slot3 := Slot{
					id: 					i+3,
					vehicleType:                CAR,
					registrationNumber:         "",
					isOccupied:                 false,
					isOccupiedByHigherPriority: false,
				}
				(*slots)[i] = slot
				(*slots)[i+1] = slot2
				(*slots)[i+2] = slot3

		}
	}

	err := errors.New("Slot id is not correct")
	return false, err
}

func NewParkingLot(vehicleSlots []VehicleType) ParkingLot  {
	p := ParkingLot{}
	slots := make([]Slot, 0)
	for i, v := range vehicleSlots {
		slot := Slot{
			id : i+1,
			vehicleType: v,
		}
		slots = append(slots, slot)
	}
	p.slots = &slots
	p.tickets = make([]Ticket, 0)
	vehicles := make([]Vehicle, 0)
	p.vehicles = vehicles

	return p
}

/*
ParkingLot struct
slots *[]Slots
tickets *[]Tickets

functions
NewParkingLot
AssignASlot
ClearASlot
DisplayStatus
CheckSlotAvailability

*/
func(p *ParkingLot) AssignASlot(vehicle Vehicle) (Slot, error) {
	slot, err := vehicle.ParkVehicle(p.slots)
	if err != nil {
		return slot, err
	}
	p.vehicles = append(p.vehicles, vehicle)
	fmt.Println("vehicle is Assigned")
	return slot, err
}

func(p *ParkingLot) ClearASlot(slotId int, vehicleRegistration string) error {
	slot := Slot{}
	if vehicleRegistration != "" {
		for _, v := range *p.slots {
			if (v.registrationNumber == vehicleRegistration) {
				slot = v
			}
		}
	}

	if slotId != 0 {
		for _, v := range *p.slots {
			if (v.id == slotId) {
				slot = v
			}
		}
	}
	
	for _, v := range p.vehicles {
		if slot.registrationNumber == v.getRegistrationNumber() {
			v.RemoveVehicle(p.slots, slot)

			fmt.Println("vehicle is cleared")
			return nil
		}
	}


	return errors.New("No SlotId or RegistrationNumber is invalid")
}

func(p *ParkingLot) DisplayStatus(vehicle Vehicle) {
    for _, v := range *p.slots {
		fmt.Printf("Slot %d \n", v.id)
		fmt.Printf("Slot Availability %d \n", v.isOccupied)
		fmt.Printf("Slot Availability %f \n", v.isOccupied)
		fmt.Printf("Slot Availability %d \n", v.isOccupiedByHigherPriority)
		fmt.Println()
	}
}

func(p *ParkingLot) CheckSlotAvailability(vehicle Vehicle) {

}

func main()  {
	vehicleTypes := []VehicleType{BUS, MOTORCYCLE, CAR, CAR, MOTORCYCLE, BUS, CAR, CAR, CAR}
	parkingLot := NewParkingLot(vehicleTypes)

	vehicleBus := Bus{}
	vehicleBus.setRegistrationNUmber("123")
	vehicleBus.setVehicleType(BUS)


	vehicleCar := Car{}
	vehicleBus.setRegistrationNUmber("12323")
	vehicleBus.setVehicleType(CAR)


	vehicleMotorCycle := Motorcycle{}
	vehicleBus.setRegistrationNUmber("12342")
	vehicleBus.setVehicleType(MOTORCYCLE)

	vehicleBus2 := Bus{}
	vehicleBus.setRegistrationNUmber("12233")
	vehicleBus.setVehicleType(BUS)

	parkingLot.AssignASlot(&vehicleBus)
	parkingLot.AssignASlot(&vehicleBus2)
	parkingLot.AssignASlot(&vehicleCar)
	parkingLot.ClearASlot(0, vehicleBus2.registrationNumber)
	parkingLot.AssignASlot(&vehicleMotorCycle)
}