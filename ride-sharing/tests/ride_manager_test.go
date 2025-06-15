package tests

// import (
// 	"ride_sharing/models"
// 	"ride_sharing/services"
// 	"testing"
// )

// func TestOfferRide(t *testing.T) {
//     rm := services.RideManager{}
//     rm.Users = append(rm.Users, models.User{
//         Name: "Rohan",
//         Vehicle: []models.Vehicle{
//             {Model: "Swift", NumberPlate: "KA-01-12345"},
//         },
//     })

//     rm.OfferRide("Rohan", "Hyderabad", "Bangalore", "Swift", 2)
//     if len(rm.Rides) != 1 {
//         t.Errorf("Expected 1 ride, got %d", len(rm.Rides))
//     }
// }

// func TestSelectRide(t *testing.T) {
//     rm := services.RideManager{}
//     driver := models.User{
//         Name: "Rohan",
//         Vehicle: []models.Vehicle{
//             {Model: "Swift", NumberPlate: "KA-01-12345"},
//         },
//     }
//     rm.Users = append(rm.Users, driver)
//     rm.Rides = append(rm.Rides, models.Ride{
//         Driver:        rm.Users[0],
//         Vehicle:       driver.Vehicle[0],
//         Source:        "Hyderabad",
//         Destination:   "Bangalore",
//         AvailableSeats: 2,
//         IsActive:      true,
//     })

//     ride := rm.SelectRide("TestUser", "Hyderabad", "Bangalore", "Most Vacant", 1)
//     if ride == nil {
//         t.Errorf("Expected to find a ride, got nil")
//     }
// }
