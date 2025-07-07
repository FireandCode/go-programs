package main

// import "fmt"

// type Car2 struct {
// 	Model     string
// 	Wheel     string
// 	Dashboard string
// }

// type ICarBuilder interface {
// 	setModel(model string) ICarBuilder
// 	setWheel(wheel string) ICarBuilder
// 	AddDashboard(dashboard string) ICarBuilder
// 	Build() Car2
// }

// type CarBuilder struct {
// 	car Car2
// }

// func NewCarBuilder() ICarBuilder {
// 	return &CarBuilder{}
// }

// func (cb *CarBuilder) setModel(model string) ICarBuilder {
// 	cb.car.Model = model
// 	return cb
// }
// func (cb *CarBuilder) setWheel(wheel string) ICarBuilder {
// 	cb.car.Wheel = wheel
// 	return cb
// }

// func (cb *CarBuilder) AddDashboard(dashboard string) ICarBuilder {
// 	cb.car.Dashboard = dashboard
// 	return cb
// }

// func (cb *CarBuilder) Build() Car2 {

// 	return cb.car
// }

// func main() {
// 	car := NewCarBuilder().setModel("BMW").setWheel("4X4").AddDashboard("prime").Build()

// 	fmt.Println(car.Model)
// }