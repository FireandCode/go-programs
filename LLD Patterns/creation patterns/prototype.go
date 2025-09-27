package main

// import "fmt"

// type IVehicle interface {
// 	getModel() string
// 	setModel(model string)
// 	clone() IVehicle
// }

// type Motor struct {
// 	model string
// }

// func (m *Motor) getModel() string {
// 	return m.model
// }

// func (m *Motor) setModel(model string) {
// 	m.model = model
// }
// func NewMotor() IVehicle {
// 	return &Motor{model: "kawasaki"}
// }

// func (m *Motor) clone() IVehicle {
// 	return &Motor{
// 		model: m.model,
// 	}
// }

// func main() {
// 	motor := NewMotor()
// 	motor1 := motor.clone()

// 	motor1.setModel("Ninja")
// 	fmt.Println(motor.getModel())
// 	fmt.Println(motor1.getModel())
// }