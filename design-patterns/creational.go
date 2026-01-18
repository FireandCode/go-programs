package main

import (
	"fmt"
	"sync"
)

//factory pattern

type ICar interface{
	SetMaxSpeed(speed int)
	GetMaxSpeed() int 
}

type BMWCar struct{
	speed int
}

func (b *BMWCar) SetMaxSpeed(speed int) {
	b.speed = speed
}

func (b *BMWCar) GetMaxSpeed() int {
	return b.speed
}

type FerraiCar struct{
	speed int
}

func (b *FerraiCar) SetMaxSpeed(speed int) {
	b.speed = speed
}

func (b *FerraiCar) GetMaxSpeed() int {
	return b.speed
}

func CarFactory(carType string) ICar {

	car, ok := carMapping[carType]
	if !ok {
		fmt.Println("Invalid Car Type")
		return nil
	}
	return car()
}

var carMapping = map[string] func() ICar {
	"bmw": func() ICar { return &BMWCar{}},
	"ferrai": func() ICar {return &FerraiCar{}},
}

//builder pattern
type House struct {
	color string
	height int
}
type HouseBuilder struct {
	house *House
}

func NewHouseBuilder() *HouseBuilder {
	return &HouseBuilder{house: &House{}}
}

func (hb *HouseBuilder) SetColor(color string) *HouseBuilder {
	hb.house.color = color
	return hb
}

func (hb *HouseBuilder) SetHeight(height int) *HouseBuilder {
	hb.house.height = height
	return hb
}

func (hb *HouseBuilder) Build() *House {
	result := hb.house
	hb.house = &House{} // auto-reset
	return result
}

//Prototype Pattern
//Prototype Pattern is a design pattern 
// where you create new objects by copying an existing object (a prototype) instead of building from scratch


type IVehicle interface{
	SetAverage(average int)
	GetAverage() int
	Clone() IVehicle
}

type Vehicle struct {
	average int
	brand string 
	color string 
}

func (v *Vehicle) SetAverage(average int) {
	v.average = average
}

func (v *Vehicle) GetAverage() int {
	return v.average
}

func(v *Vehicle) Clone() IVehicle{
	newV := *v 
	return &newV
}


func NewVehicle() IVehicle {
	return &Vehicle{
		color: "blue",
		brand: "lambourgini",
		average: 100,
	}

}

//singelton pattern
//
type DB struct {
	connection int64 
}

var (
	once sync.Once
	db *DB
	count int64
)


func NewInstance() *DB {
	count = count+1
		once.Do(
		func () {
			db = &DB{
				connection: count,
			}
	})
	fmt.Println(db.connection)
	return db
}


// func main() {
	// go NewDB()
	// go NewDB()
	// time.Sleep(1*time.Second)
	// bmw1 := CarFactory("bmw")
	// bmw1.SetMaxSpeed(111)
	// bmw2 := CarFactory("bmw")
	// fmt.Println(bmw1.GetMaxSpeed())
	// fmt.Println(bmw2.GetMaxSpeed())
	// fmt.Println("HELLO WORLD")


	// hb := NewHouseBuilder()

	// house1 := hb.SetColor("blue").SetHeight(150).Build()
	// house2 := hb.SetColor("pink").SetHeight(119).Build()

	// fmt.Println(house1.color)
	// fmt.Println(house2.color)

	// vehicle1 := NewVehicle()
	
	// vehicle2 := vehicle1.Clone()
	// vehicle2.SetAverage(111)
	// fmt.Println(vehicle1.GetAverage())
	// fmt.Println(vehicle2.GetAverage())
// }