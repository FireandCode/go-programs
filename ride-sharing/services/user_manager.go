package services

import (
	"ride_sharing/models"
)

// UserManager handles user-related operations
type UserManager struct {
	Users []*models.User // Slice of pointers, as user data might be shared and modified
}

// AddUser adds a new user to the system
func (um *UserManager) AddUser(name, gender string, age int) {
	um.Users = append(um.Users, &models.User{
		Name:   name,
		Gender: gender,
		Age:    age,
	})
}

// AddVehicle adds a vehicle to an existing user
func (um *UserManager) AddVehicle(userName, model, numberPlate string) bool {
	for _, user := range um.Users {
		if user.Name == userName {
			user.Vehicle = append(user.Vehicle, &models.Vehicle{
				Model:       model,
				NumberPlate: numberPlate,
			})
			return true
		}
	}
	return false
}
