package services

import (
	"fmt"
	"ride_sharing/models"
)

// RideManager handles ride-related operations
type RideManager struct {
	Rides []*models.Ride    // Slice of pointers, as rides might be modified or queried
	Users []*models.User    // Shared reference to UserManager's users
}

// OfferRide allows a user to offer a ride
func (rm *RideManager) OfferRide(driverName, vehicleModel, origin, destination string, seats int) {
	fmt.Println("rides ", len(rm.Users))
	for _, user := range rm.Users {
		if user.Name == driverName {
			fmt.Println("rides ", len(rm.Rides))
			for _, v := range user.Vehicle {
				if v.Model == vehicleModel {
					fmt.Println("rides ", len(rm.Rides))
					for _, r := range rm.Rides {
						if r.Driver.Name == driverName && r.Vehicle.Model == vehicleModel && r.IsActive {
							fmt.Println("Ride already active for this vehicle.")
							return
						}
					}
					ride := &models.Ride{
						Driver:         user,
						Vehicle:        v,
						Origin:         origin,
						Destination:    destination,
						AvailableSeats: seats,
						IsActive:       true,
					}
					rm.Rides = append(rm.Rides, ride)
					fmt.Println("rides ", len(rm.Rides))
					user.RidesOffered++
					return
				}
			}
		}
	}
	fmt.Println("Driver or Vehicle not found.")
}

// SelectRide allows a passenger to select a ride
func (rm *RideManager) SelectRide(passengerName, origin, destination string, seats int, preferredVehicle string) {
	for _, ride := range rm.Rides {
		if ride.Origin == origin && ride.Destination == destination && ride.IsActive && ride.AvailableSeats >= seats {
			if preferredVehicle == "" || ride.Vehicle.Model == preferredVehicle {
				ride.AvailableSeats -= seats
				fmt.Printf("%s successfully booked the ride.\n", passengerName)
				return
			}
		}
	}
	fmt.Println("No suitable ride found.")
}
