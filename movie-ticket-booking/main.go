package main

import (
	"errors"
	"fmt"
	"time"
)

/*
Movie - struct
- Name
- duration
- language

Theater - struct
- []Screens
- Location
- Name
- Id

Seat - struct
- Id
- Type
- isAvailable

Pricing - struct
- seatType
- price

Screen - struct
- Movie
- []Timing
- []Seats
- Id

Timing - struct
StartTime int
EndTime  int

User - struct
- Id
- []Ticket

PaymentMethod - Enum

Ticket - struct
- Id
- PaymentMethod
- amount
- userId
- MovieName
- SeatId
- ScreenId
- TheaterId

BookmyShow - struct
- []Theater
- []Ticket

Functions
- NewBookMyShow
- AddATheater
- AddAScreen
- AddAMovie
- BookATicket
- CancelTicket
- ModifyTicket
- ShowAvailability
- showBookingStatus
- SearchAMovie

*/
var theaterId = 0
var screenId = 0
func NewTheater(name string, location string, noOfScreens int) Theater {
	theaterId = theaterId +1
	return Theater{
		id: theaterId,
		name: name,
		location: location,
		screens: make([]Screen, noOfScreens),
	}
}

func NewScreen(seats map[SeatType]int) Screen {
	screenId = screenId +1
	finalSeats := make(map[SeatType][]Seat)

	finalSeats[SILVER] = make([]Seat, seats[SILVER])
	finalSeats[GOLD] = make([]Seat, seats[GOLD])
	finalSeats[PREMIUM] = make([]Seat, seats[PREMIUM])
	return Screen{
		id: screenId,
		seats: finalSeats,
	}

}

func NewBookMyShow(theaters []Theater) (BookmyShow)  {
	bookMyShow := BookmyShow{}

	bookMyShow.theaters = theaters
	bookMyShow.tickets = make([]Ticket, 0)

	return bookMyShow
}

func(b *BookmyShow) AddATheater(name string, location string, noOfScreens int)  {
	theater := NewTheater(name, location, noOfScreens)

	b.theaters = append(b.theaters, theater)

}

func(b *BookmyShow) AddAScreen(theaterId int, seats map[SeatType]int)  {
	screen := NewScreen(seats)

	for i, v := range b.theaters {
		if v.id == theaterId {
			theater := v
			theater.screens = append(theater.screens, screen)
			(b.theaters)[i] = theater
		}
	}
}

func NewMovie(name string, duration int) Movie {
	return Movie{
		name: name,
		duration: duration,
	}
}

func(b *BookmyShow) AddAMovie(name string, screenId int, theaterId int, duration int, timings []Timing)  {
	movie := NewMovie(name, duration)

	for i, v := range b.theaters {
		if v.id == theaterId {
			for j, s := range v.screens {
				if s.id == screenId {
					sNo := make(map[SeatType]int, 0)
					sNo[GOLD] = len(s.seats[GOLD])
					sNo[SILVER] = len(s.seats[SILVER])
					sNo[PREMIUM] = len(s.seats[PREMIUM])
					screen := NewScreen(sNo)

					screen.movie = movie
					screen.id = s.id
					screen.timings = timings

					(b.theaters)[i].screens[j] = screen
				}
			}
		}
	}
}

/*
- BookATicket
- CancelTicket
- ModifyTicket
- ShowAvailability
- showBookingStatus
- SearchAMovie
*/

func (b *BookmyShow) BookATicket(userId int, theaterId int, screenId int, movieName string, seatType SeatType) (Ticket, error) {
    for i, theater := range b.theaters {
        if theater.id == theaterId {
            for j, screen := range theater.screens {
                if screen.id == screenId && screen.movie.name == movieName {
                    for k, seat := range screen.seats[seatType] {
                        if seat.isAvailable {
                            b.theaters[i].screens[j].seats[seatType][k].isAvailable = false
                            ticket := Ticket{
                                id:           len(b.tickets) + 1,
                                paymentMethod: UPI,  // Default payment method (can be modified later)
                                amount:        getPrice(seatType),
                                userId:        userId,
                                MovieName:     movieName,
                                seatId:        seat.id,
                                screenId:      screenId,
                                theaterId:     theaterId,
                                date:          int(time.Now().Unix()),
                            }
                            b.tickets = append(b.tickets, ticket)
                            return ticket, nil
                        }
                    }
                    return Ticket{}, errors.New("No available seats of the selected type")
                }
            }
        }
    }
    return Ticket{}, errors.New("Invalid booking request")
}

func (b *BookmyShow) CancelTicket(ticketId int) error {
    for i, ticket := range b.tickets {
        if ticket.id == ticketId {
            for j, theater := range b.theaters {
                if theater.id == ticket.theaterId {
                    for k, screen := range theater.screens {
                        if screen.id == ticket.screenId {
                            for l, seat := range screen.seats {
                                if seat[l].id == ticket.seatId {
                                    b.theaters[j].screens[k].seats[l][seat[l].id].isAvailable = true
                                    b.tickets = append(b.tickets[:i], b.tickets[i+1:]...)
                                    return nil
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    return errors.New("Ticket not found")
}


func (b *BookmyShow) ModifyTicket(ticketId int, newSeatType SeatType) error {
    for _, ticket := range b.tickets {
        if ticket.id == ticketId {
            err := b.CancelTicket(ticketId)
            if err != nil {
                return err
            }
            _, err = b.BookATicket(ticket.userId, ticket.theaterId, ticket.screenId, ticket.MovieName, newSeatType)
            return err
        }
    }
    return errors.New("Ticket not found")
}

func (b *BookmyShow) ShowAvailability(theaterId int, screenId int) map[SeatType]int {
    availability := make(map[SeatType]int)
    for _, theater := range b.theaters {
        if theater.id == theaterId {
            for _, screen := range theater.screens {
                if screen.id == screenId {
                    for seatType, seats := range screen.seats {
                        count := 0
                        for _, seat := range seats {
                            if seat.isAvailable {
                                count++
                            }
                        }
                        availability[seatType] = count
                    }
                }
            }
        }
    }
    return availability
}

func (b *BookmyShow) ShowBookingStatus() []Ticket {
    return b.tickets
}

func (b *BookmyShow) SearchAMovie(movieName string) []Theater {
    var result []Theater
    for _, theater := range b.theaters {
        for _, screen := range theater.screens {
            if screen.movie.name == movieName {
                result = append(result, theater)
            }
        }
    }
    return result
}

func getPrice(seatType SeatType) int {
    switch seatType {
    case SILVER:
        return 100
    case GOLD:
        return 200
    case PREMIUM:
        return 300
    default:
        return 0
    }
}

type BookmyShow struct {
	theaters []Theater
	tickets []Ticket
}

type Theater struct {
	id int
	location string
	name string
	screens []Screen
}

type Movie struct {
	name string
	duration int
}

type Timing struct {
	startTime int 
	endTime int 
}
type Screen struct {
	id int 
	movie Movie
	timings []Timing
	seats map[SeatType][]Seat
}

type User struct {
	id int
	tickets []Ticket
}

type SeatType int 

const (
	SILVER SeatType = iota
	GOLD 
	PREMIUM
)

type Pricing struct {
	seatType SeatType
	price int
}

type Seat struct {
	id int 
	seatType SeatType
	isAvailable bool 
}

type PaymentMethod int 

const (
	UPI PaymentMethod = iota
	CARD 
	NETBANKING
)

type Ticket struct {
	id int
	paymentMethod PaymentMethod
	amount int
	userId int
	MovieName string
	seatId int
	screenId int
	theaterId int
	date int
}

func main() {
    bms := BookmyShow{}
    
    theater := Theater{id: 1, name: "IMAX", location: "NYC", screens: []Screen{}}
    bms.theaters = append(bms.theaters, theater)
    
    screenSeats := map[SeatType]int{SILVER: 50, GOLD: 30, PREMIUM: 20}
    bms.AddAScreen(1, screenSeats)
    
    timings := []Timing{{startTime: 1800, endTime: 2000}}
    bms.AddAMovie("Inception", 1, 1, 120, timings)
    
    ticket, err := bms.BookATicket(1, 1, 1, "Inception", GOLD)
    if err != nil {
        fmt.Println("Booking failed:", err)
    } else {
        fmt.Println("Ticket booked:", ticket)
    }
    
    fmt.Println("Current Booking Status:", bms.ShowBookingStatus())
}