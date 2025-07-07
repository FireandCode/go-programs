package main

// import "fmt"

// type IVehicle interface {
// 	getModel() string
// }

// type Vehicle struct {
// 	model string
// }

// func (v *Vehicle) getModel() string {
// 	return v.model
// }

// type Car struct {
// 	Vehicle
// }

// type Truck struct {
// 	Vehicle
// }

// func NewCar() IVehicle {
// 	return &Car{
// 		Vehicle: Vehicle{
// 			model: "BMW",
// 		},
// 	}
// }

// func NewTruck() IVehicle {
// 	return &Truck{
// 		Vehicle: Vehicle{
// 			model: "Essay",
// 		},
// 	}
// }
// func NewVehicle(vehicleType string) IVehicle {
// 	if vehicleType == "car" {
// 		return NewCar()
// 	}
// 	if vehicleType == "truck" {
// 		return NewTruck()
// 	}

// 	return nil
// }

// func main() {
// 	vehicle1 := NewVehicle("car")
// 	vehicle2 := NewVehicle("truck")

// 	fmt.Println(vehicle1.getModel())
// 	fmt.Println(vehicle2.getModel())

// }