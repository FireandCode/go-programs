package main

import (
	"fmt"
	"strconv"
)

//Adapter Pattern
type IOldPrinter interface {
	Print(pages int)
}
type OldPrinter struct {
	color string 
	cost int 
}

func(op *OldPrinter) Print(pages int) {
	fmt.Println("cost ", op.cost*pages)
	fmt.Println("coloring scheme ", op.color)
	fmt.Println("no. of pages", pages)
}

type IPrinter interface {
	Print(pages string)
}

type PrinterAdapter struct {
	op IOldPrinter
}

func(pa *PrinterAdapter) Print(pages string) {
	pagesInt, err := strconv.ParseInt(pages, 10, 64)
	if err != nil {
		fmt.Println(err.Error())
		return 
	}
	pa.op.Print(int(pagesInt))
}

func NewOldPrinter() IOldPrinter {
	return &OldPrinter{
		color: "black/white",
		cost: 22,
	}
}

func NewPrinterAdapter(op IOldPrinter) IPrinter {
	return &PrinterAdapter{
		op: op,
	}
}

/*
Decorator Pattern
we wrap an object around another object to increase it's capabilities.

*/

type ICoffee interface{
	Cost() int 
}

type BaseCoffee struct {
}

func(bc *BaseCoffee) Cost() int {
	return 100
}

type Milk struct {
	c ICoffee
}

func(m *Milk) Cost() int {
	return m.c.Cost() + 10
} 

type Sugar struct {
	c ICoffee
}

func(s *Sugar) Cost() int {
	return s.c.Cost() + 15
}

func NewCoffee() ICoffee {
	return &BaseCoffee{}
}

// func main() {

// 	coffee := NewCoffee()
// 	coffee = &Milk{c: coffee}
// 	coffee = &Sugar{c: coffee}

// 	fmt.Println(coffee.Cost())
	
// 	client := "222"
// 	op := NewOldPrinter()

// 	pAdapter := NewPrinterAdapter(op)

// 	pAdapter.Print(client)

// }