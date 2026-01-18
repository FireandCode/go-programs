package main

import "fmt"

// Strategy interface
type IPriceAlgo interface {
	Calculate() int
}

// Place strategy
type PlaceAlgo struct {
	place string
	cost  int
}

func (pa *PlaceAlgo) Calculate() int {
	switch pa.place {
	case "boho":
		return pa.cost * 10
	default:
		return pa.cost * 4
	}
}

// Drink strategy
type DrinkAlgo struct {
	drinks int
	cost   int
}

func (da *DrinkAlgo) Calculate() int {
	return da.cost * da.drinks
}

// Strategy Registry that accepts dynamic params
var algos = map[string]func(input map[string]any) IPriceAlgo{
	"place": func(input map[string]any) IPriceAlgo {
		return &PlaceAlgo{
			place: input["place"].(string),
			cost:  input["cost"].(int),
		}
	},
	"drink": func(input map[string]any) IPriceAlgo {
		return &DrinkAlgo{
			drinks: input["drinks"].(int),
			cost:   input["cost"].(int),
		}
	},
}

// Context
type PriceCalculator struct {
	algo IPriceAlgo
}

func NewPriceCalculator(algo IPriceAlgo) *PriceCalculator {
	return &PriceCalculator{algo: algo}
}

func (pc *PriceCalculator) ComputePrice() int {
	return pc.algo.Calculate()
}


//Observer Pattern
//Shop -> products -> 
//subscribe to notification for new products
//adding a new product to shop then we will need to notify the subscribed users

//User ->

type ProductEvent struct {
    ID string
}

type ISubscriber interface {
    GetID() int64
    Notify(event ProductEvent) bool
}

type Shop struct {
    products    []string
    subscribers map[int64]ISubscriber
}

func NewShop() *Shop {
    return &Shop{
        products:    []string{},
        subscribers: make(map[int64]ISubscriber),
    }
}

func (s *Shop) Publish(product string) {
    s.products = append(s.products, product)
    event := ProductEvent{ID: product}

    for _, sub := range s.subscribers {
        sub.Notify(event)
    }
}

func (s *Shop) Register(sub ISubscriber) {
    s.subscribers[sub.GetID()] = sub
}

func (s *Shop) Remove(sub ISubscriber) {
    delete(s.subscribers, sub.GetID())
}

// subscribers
type BaseSubscriber struct{ Id int64 }
func (bs *BaseSubscriber) GetID() int64 { return bs.Id }

type Person struct{ BaseSubscriber }
func (p *Person) Notify(event ProductEvent) bool {
    fmt.Println("Person notified of", event.ID)
    return true
}

type Company struct{ BaseSubscriber }
func (c *Company) Notify(event ProductEvent) bool {
    fmt.Println("Company bulk-buy alert:", event.ID)
    return true
}



func main() {

//Observor Pattern
	shop := Shop{
		products: []string{"coffee"},
		subscribers: make(map[int64]ISubscriber),
	}
	person := &Person{
		BaseSubscriber: BaseSubscriber{
			Id: 1,
		},
	}
	company := &Company{
		BaseSubscriber{
			Id: 2,
		},
	}
	shop.Register(person)
	shop.Register(company)

	shop.Publish("ice-cream")


	// //---------------------------------------
	// // Pick strategy dynamically
	// clientChoice := "drink" // read from REST API, CLI, config, etc.
	// //---------------------------------------

	// //---------------------------------------
	// // Provide dynamic params
	// params := map[string]any{
	// 	"drinks": 3,
	// 	"cost":   15,
	// }
	// //---------------------------------------

	// algoFactory, ok := algos[clientChoice]
	// if !ok {
	// 	panic("strategy not found")
	// }

	// chosenAlgo := algoFactory(params)
	// calculator := NewPriceCalculator(chosenAlgo)

	// fmt.Println("Final cost:", calculator.ComputePrice())
}
