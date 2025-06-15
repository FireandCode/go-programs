package main

import (
	"fmt"
	"ride_sharing/services"
)

func main() {
	// Initialize managers
	userManager := &services.UserManager{}
	rideManager := services.RideManager{
		Users: userManager.Users,
	}

	// Add Users
	userManager.AddUser("Rahul", "Male", 30)
	userManager.AddUser("Nandini", "Female", 25)

	// Add Vehicles
	userManager.AddVehicle("Rahul", "Swift", "KA-01-12345")
	userManager.AddVehicle("Nandini", "Polo", "KA-02-54321")

	// Offer Rides
	rideManager.OfferRide("Rahul", "Swift", "Bangalore", "Mysore", 3)
	rideManager.OfferRide("Nandini", "Polo", "Bangalore", "Chennai", 2)

	// Select Rides
	rideManager.SelectRide("Rahul", "Bangalore", "Chennai", 1, "Polo")

	// Print Ride Statistics
	fmt.Println("Total Users:", len(userManager.Users))
	fmt.Println("Total Rides Offered:", len(rideManager.Rides))
}
