package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"
)

/*

🔍 Problem Statement: Design a Food Ordering Backend
Design the backend system for a simplified Food Booking Application. The platform enables users to discover restaurants, explore their menus, and place food orders from available options.
💡 Core Functional Expectations:
Restaurant Discovery:
Users should be able to search for restaurants by Name. (should easily support other filters also).
Only nearby restaurants (based on the user's current location) should be returned in search results.


Ordering Flow:
Users can explore a restaurant’s menu and choose food items to prepare for an order.
The system should support the ability to place an order with a selected restaurant.
Note: An order must be associated with a single restaurant.


Capacity Constraints:
Each restaurant has a limit on how many orders it can process concurrently.
If the number of active orders exceeds this threshold, the restaurant should not accept new orders until some of the current orders are marked as completed.

Flows
1. static menu for restaurants

1. onboarding a restaurant
	-> name, location(x, y), menu, activeOrders
	FoodDelivery(Admin) -> newRestaurant()
2. onboarding a user
	-> client -> newUser(User{})
3. search for all the restaurants
	-> client -> searchRestaurant(location, name)
		filter - name, // should be extensible
		Filters - name
	 -> name , id
4. checking the menu using id
	-> //pass the id
	-> menu

5. user placing an order
	-> //pass -> id, list of items{}
	-> response -> order_id, total_amount, items, status, restaurant_id

struct

FoodDelivery()
resturants map[string]Resturant
users 	map[string]User
orders map[string]Order

Restaurant
Id
name
menu []Items
activeOrders int
location

*/

type FoodDelivery struct {
	resturants map[string]Resturant
	users map[int]User
	orders map[int]Order

}

type Resturant struct {
	id int 
	name string
	location int
	menu map[int]MenuItem
	activeOrders int
	threshold int
	mu *sync.Mutex
}

type MenuItem struct {
	id int
	name string
	price int
	description string
}

type User struct {
	id int
	name string
	location int
}

type Item struct {
	menuID int 
	count int
}

type Order struct {
	id int
	userId int
	resturantID int
	items []Item
	totalAmount int
	status string 
}

const (
	LOCATION_THRESHOLD = 5
)

type RestaurantFilter interface {
    Apply(r Resturant, userLocation int) bool
}

type NameFilter struct {
    Name string
}

func (nf NameFilter) Apply(r Resturant, _ int) bool {
    return r.name == nf.Name
}

type LocationFilter struct {}

func (lf LocationFilter) Apply(r Resturant, userLocation int) bool {
    diff := r.location - userLocation
    if diff < 0 {
        diff = -diff
    }
    return diff <= LOCATION_THRESHOLD
}

func NewFoodDelivery() FoodDelivery {
	foodDelivery := FoodDelivery{
		resturants: make(map[string]Resturant),
		users: make(map[int]User),
		orders: make(map[int]Order),
	}	
	return foodDelivery
}

func generateRandomID() int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(1000)))
	return int(n.Int64())
}

func NewResturant(name string, menu []MenuItem, location int) Resturant {
	 menuR := make(map[int]MenuItem)
	 for i := 0; i < len(menu) ; i++ {
		menuR[menu[i].id] = menu[i]
	 }
	
	 restaurant := Resturant{
		id : generateRandomID(),
		name: name,
		menu: menuR,
		location: location,
		activeOrders: 0,
		threshold: 2,
		mu: &sync.Mutex{},
	 }
	return restaurant
}

 func NewUser(name string, location int) User {
	return User{
		id: generateRandomID(),
		name: name,
		location: location,
	}

 }


func (fd *FoodDelivery) SearchRestaurants(userLocation int, filters ...RestaurantFilter) []Resturant {
    results := []Resturant{}

    for _, restaurant := range fd.resturants {
        matches := true
        for _, filter := range filters {
            if !filter.Apply(restaurant, userLocation) {
                matches = false
                break
            }
        }
        if matches {
            results = append(results, restaurant)
        }
    }

    return results
}


func (fd *FoodDelivery) GetResturantByID(id int  ) Resturant {
	resturant := Resturant{}
	for _, v := range fd.resturants {
		if v.id == id{
			resturant = v
		}
	}

	return resturant
}

func verifyMenu(restaurantMenu map[int]MenuItem, userMenu []Item) error {
    for _, item := range userMenu {
        if _, exists := restaurantMenu[item.menuID]; !exists {
            return fmt.Errorf("menu item %d not found", item.menuID)
        }
    }
    return nil
}

func calculateTotalAmount(resurantMenu  map[int]MenuItem, userMenu []Item) int {
	totalAmount := 0
	for i := 0; i < len(userMenu); i++ {
		item := resurantMenu[userMenu[i].menuID]
		totalAmount += item.price*userMenu[i].count
	}
	return totalAmount
}

func(fd *FoodDelivery) PlaceAOrder(restaurantID int, userId int, userMenu []Item) (Order, error)  {
	
	restaurant := fd.GetResturantByID(restaurantID)
	if restaurant.id == 0 {
		return Order{}, errors.New("no resturant found with id")
	}

	if fd.users[userId].id == 0 {
		return Order{}, errors.New("no user found with id")
	}
	
	//mutex
	restaurant.mu.Lock()
	defer restaurant.mu.Unlock()

	if fd.resturants[restaurant.name].activeOrders >= restaurant.threshold {
		return Order{} , errors.New("resturant is not taking any more order at this point")
	}

	 err := verifyMenu(restaurant.menu, userMenu)
	if err != nil {
		return Order{}, err
	}

	order := Order{
		id: generateRandomID(),
		resturantID: restaurantID,
		userId: userId,
		items: userMenu,
		totalAmount: calculateTotalAmount(restaurant.menu, userMenu),
		status: "initated",
	}
	
	fd.orders[order.id] = order
	restaurant= fd.resturants[restaurant.name]
	restaurant.activeOrders++

	fd.resturants[restaurant.name] = restaurant
	return order, nil
}

func main() {
	foodDelivery := NewFoodDelivery()
	menu := []MenuItem{
	{
		id: 1,
		name: "chicker burger",
		price: 100,
	},
	{
		id: 2,
		name: "veg burger",
		price: 80,
	},
	}
	restaurant := NewResturant("burgers", menu, 3)
	foodDelivery.resturants[restaurant.name] = restaurant
	user := NewUser("rahul", 5)
	foodDelivery.users[user.id] = user

	item := []Item{
		{
			menuID: 1,
			count: 2,
		},
	}
	order1, err := foodDelivery.PlaceAOrder(restaurant.id, user.id, item)
	if err != nil{
		fmt.Println(err.Error())
	}

		order2, err := foodDelivery.PlaceAOrder(restaurant.id, user.id, item)
	if err != nil{
		fmt.Println(err.Error())
	}


		order3, err := foodDelivery.PlaceAOrder(restaurant.id, user.id, item)
	if err != nil{
		fmt.Println(err.Error())
	}

	fmt.Println(order1.id, order2.id, order3.id)

	searchResults := foodDelivery.SearchRestaurants(
    user.location,
    NameFilter{Name: "burgers"},
    LocationFilter{},
)

for _, r := range searchResults {
    fmt.Printf("Found Restaurant: %s (ID: %d)\n", r.name, r.id)
}

}
