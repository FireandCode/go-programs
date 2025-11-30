package main

import (
	"errors"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"time"
)

/*
User
Ticket

*/

type User struct {
	id string 
	name string
	tickets map[string]*Ticket
	mut sync.Mutex
}

type Ticket struct {
	id string
	flightID string
	userID string
	seatID int
	cost int
}

type Flight struct {
	id string
	seats []*Seat
	source string 
	destination string 
	price int
}

type Seat struct {
	id int 
	passenger Passenger
	status string 
	mut sync.Mutex
}

type Passenger struct {
	name string 
	age int
	identity string 
}

type AirlineBookingSystem struct {
	flights map[string]Flight
	users map[string]*User
}

type Filter interface{
	Apply(flights []Flight) []Flight
}

type ColumnFilter struct {
	column string
	value string
}

func (cf *ColumnFilter) Apply(flights []Flight) []Flight {
	filtered := make([]Flight, 0)
	for _, flight := range flights {
		v := reflect.ValueOf(flight)
		field := v.FieldByName(cf.column)
		if field.IsValid() && field.Interface() == cf.value {
			filtered = append(filtered, flight)
		}
	}
	return filtered
}

type FlightComparator func(a, b Flight) bool

var sortComparators = map[string]map[string]FlightComparator{
	"Price": {
		"asc":  func(a, b Flight) bool { return a.price < b.price },
		"desc": func(a, b Flight) bool { return a.price > b.price },
	},
	"Source": {
		"asc":  func(a, b Flight) bool { return a.source < b.source },
		"desc": func(a, b Flight) bool { return a.source > b.source },
	},
	// Add more fields here without touching SortFilter.Apply
}

type SortFilter struct {
	basedOn string 
	sort string 
}

func (sf *SortFilter) Apply(flights []Flight) []Flight {
	if fieldSorters, ok := sortComparators[sf.basedOn]; ok {
		if comparator, ok := fieldSorters[sf.sort]; ok {
			sort.Slice(flights, func(i, j int) bool {
				return comparator(flights[i], flights[j])
			})
		}
	}
	return flights
}

type Status int 

const (
	Available Status = iota
	Booked 
	NotAvailable
)

/*
user flows -> search a flight
*/

func(ab *AirlineBookingSystem) SearchFlight(filters ...Filter) []Flight {
	flights := make([]Flight, 0)
	for _, v := range ab.flights {
		flights = append(flights, v)
	}
	for _, v := range filters {
		flights = v.Apply(flights)
	}
	return flights
}

func(ab *AirlineBookingSystem) BookFlight(flightID string, userID string, passenger Passenger, seatID int) (Ticket, error){
	flight := ab.flights[flightID]
	user := ab.users[userID]
	seat := flight.seats[seatID]
	seat.mut.Lock()
	defer seat.mut.Unlock()
	if(seat.status != "available") {
		return Ticket{}, errors.New("seat is already booked")
	}
	seat.status = "booked"
	seat.passenger = passenger
	ticket := Ticket{
		id: generateRandomID(),
		userID: userID,
		flightID: flightID,
		seatID: seatID,
		cost: flight.price,
	}
	user.tickets[ticket.id] = &ticket

	return ticket, nil
}

func NewAirlineBookingSystem() AirlineBookingSystem {
	return AirlineBookingSystem{
		flights: make(map[string]Flight),
		users:  make(map[string]*User),
	}
}

func generateRandomID() string {
	return strconv.FormatInt(rand.Int63n(10000*100), 10) + strconv.FormatInt(time.Now().UnixMilli(), 10)
}

func NewUser(name string) User {
	return User{
		id : generateRandomID(),
		name : name,
		tickets: make(map[string]*Ticket),
		mut: sync.Mutex{},
	}
}

func NewFlight(source string, des string, noOfSeats int, price int ) Flight {
	return Flight{
		id: generateRandomID(),
		source: source,
		destination: des,
		seats: NewSeats(noOfSeats),
		price: price,
	}
}

func NewSeat(id int) *Seat {
	return &Seat{
		id : id,
		passenger: Passenger{},
		status: "available",
		mut: sync.Mutex{},
	}
}

func NewSeats(noOfSeats int ) []*Seat {
	seats := make([]*Seat, 0)
	for i := 0; i < noOfSeats; i++ {
		seats = append(seats, NewSeat(i))
	}
	return seats
}

