package main

import "fmt"

// Target interface your code expects
type Printer interface {
	Print(msg string)
}

// Adaptee - legacy or 3rd party
type OldPrinter struct{}

func (o *OldPrinter) PrintLegacy(msg string) {
	fmt.Println("OldPrinter:", msg)
}

// Adapter - makes OldPrinter compatible with Printer
type PrinterAdapter struct {
	LegacyPrinter *OldPrinter
}

func (a *PrinterAdapter) Print(msg string) {
	// Convert call
	a.LegacyPrinter.PrintLegacy(msg)
}


func main() {
	old := &OldPrinter{}
	adapter := &PrinterAdapter{LegacyPrinter: old}

	var printer Printer = adapter
	printer.Print("Hello via adapter!")
}
