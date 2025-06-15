package main

import (
	"container/list"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*

Algo - interface{}
initalize
evaluate
initalizeGoRoutine

type BucketAlgo struct
	- bucketSize
	- tokenRate
	- numberOfToken
	-




type SlidingAlgo struct
	- duration

Route - struct
id
name
endpoint
Method

type Request struct
- id
- routeId
- userId
- isThrottled
- time time.Time


User - struct
	- id
	- map[route]limit
	- map[route]requestLeft
	- mutex
RateLimiter - struct
- Algo
- Users
- mutex
- requests
- Routes


AddARoute()

*/

const STATIC_LIMIT = 3

func (rl *RateLimiter) AddARoute(route Route) Route {
	route.id = generateRandomNumber()

	fmt.Println("new route is added ", route.id)
	rl.routes[route.id] = route

	return route
}

func (rl *RateLimiter) ChangeTheLimit(userId int, routeId int, limit int) {
	rl.users[userId].limits[routeId] = limit
	fmt.Println("Limit for userId ", userId, " is changed to ", limit, " for routeId", routeId)
}

func (rl *RateLimiter) AddAUser(user User) *User {
	user.id = generateRandomNumber()
	user.limits = make(map[int]int)
	user.requestLeft = make(map[int]int)
	user.mu = &sync.Mutex{}
	for k := range rl.routes {
		if _, exists := user.limits[k]; !exists {
			user.limits[k] = STATIC_LIMIT
		}
		user.requestLeft[k] = user.limits[k]
	}
	fmt.Println("new User is added ", user.id)
	rl.users[user.id] = &user

	return &user
}

func (rl *RateLimiter) UserRequest(userId int, routeId int) error {
	user, exists := rl.users[userId]
	if !exists {
		return errors.New("userId is invalid")
	}

	_, exists = rl.routes[routeId]
	if !exists {
		return errors.New("routeId is invalid")
	}

	if rl.algo.Evaluate(user, routeId) == false {
		fmt.Println("request is throttled")
		return errors.New("Request is throttled, please try again after some time")
	}

	//forward the request
	return nil
}

type RateLimiter struct {
	algo   Algo
	users  map[int]*User
	routes map[int]Route
	mu     *sync.Mutex
}

func NewRateLimiter(algo Algo) *RateLimiter {
	rl := RateLimiter{}
	rl.algo = algo
	rl.users = make(map[int]*User)
	rl.routes = make(map[int]Route)
	rl.mu = &sync.Mutex{}

	go rl.intializeTheLimitUpdate()

	return &rl
}

func (rl *RateLimiter) intializeTheLimitUpdate() {

	for {
		rl.algo.updateTheLimits(rl.users)
		time.Sleep(rl.algo.getDuration())
	}

}

type Queue struct {
	items *list.List
}

func NewQueue() *Queue {
	return &Queue{items: list.New()}
}

// Enqueue (Push to the back)
func (q *Queue) Enqueue(value Request) {
	q.items.PushBack(value)
}

// Dequeue (Pop from the front)
func (q *Queue) Dequeue() (Request, bool) {
	if q.items.Len() == 0 {
		return Request{}, false
	}
	front := q.items.Front()
	q.items.Remove(front)
	return front.Value.(Request), true
}

type Request struct {
	id          int
	routeId     int
	userId      int
	requestTime time.Time
	isThrottled bool
}

type Route struct {
	id       int
	name     string
	endpoint string
	method   string
}

type Algo interface {
	Initalize()
	Evaluate(user *User, routeId int) bool
	getDuration() time.Duration
	updateTheLimits(users map[int]*User)
}

type BucketAlgo struct {
	duration      time.Duration
	numberOfToken int
}

func (ba *BucketAlgo) Initalize(tokenDuration time.Duration, numberOfToken int) {
	ba.numberOfToken = numberOfToken
	ba.duration = tokenDuration
}

func (ba *BucketAlgo) Evaluate(user *User, routeId int) bool {
	user.mu.Lock()
	defer user.mu.Unlock()

	_, exists := user.limits[routeId]
	if !exists {
		user.requestLeft[routeId] = STATIC_LIMIT
		user.limits[routeId] = STATIC_LIMIT
	}

	if user.requestLeft[routeId] > 0 {
		user.requestLeft[routeId]--
		fmt.Println("Forward the Request ", routeId, " for User ", user.id)
		return true
	}
	return false
}

func (ba *BucketAlgo) getDuration() time.Duration {
	return ba.duration
}

func (ba *BucketAlgo) updateTheLimits(users map[int]*User) {
	fmt.Println("updating the token for users")
	for k := range users {
		for kk := range users[k].requestLeft {
			users[k].requestLeft[kk] = min(users[k].limits[kk], users[k].requestLeft[kk]+ba.numberOfToken)
		}
	}
	fmt.Println("updated the token for users")
}

type SlidingAlgo struct {
	duration    time.Duration
	window      time.Duration
	requestTime *Queue
}

func (sa *SlidingAlgo) Initalize() {
	sa.duration = 1*time.Second
	sa.window = 1* time.Second
	sa.requestTime = NewQueue()
}

func (sa *SlidingAlgo) getDuration() time.Duration {
	return sa.duration
}

func (sa *SlidingAlgo) Evaluate(user *User, routeId int) bool {
	user.mu.Lock()
	defer user.mu.Unlock()

	request := Request{
		id:          generateRandomNumber(),
		userId:      user.id,
		routeId:     routeId,
		requestTime: time.Now(),
		isThrottled: true,
	}

	_, exists := user.limits[routeId]
	if !exists {
		user.requestLeft[routeId] = STATIC_LIMIT
		user.limits[routeId] = STATIC_LIMIT
	}

	if user.requestLeft[routeId] > 0 {
		request.isThrottled = false
		sa.requestTime.Enqueue(request)
		user.requestLeft[routeId]--

		fmt.Println("Forward the Request ", routeId, " for User ", user.id)
		return true
	}

	sa.requestTime.Enqueue(request)
	return false
}

func (sa *SlidingAlgo) updateTheLimits(users map[int]*User) {
	windowTime := time.Now().Add(-sa.window)

	for {
		if sa.requestTime.items.Len() == 0 {
			break
		}
		request, ok := sa.requestTime.items.Front().Value.(Request)
		if !ok {
			fmt.Println("unable to fetch the value from queue")
			break
		}
		if request.requestTime.After(windowTime) {
			break
		}
		request, _ = sa.requestTime.Dequeue()

		limit, exists := users[request.userId].limits[request.routeId]
		if !exists {
			users[request.userId].limits[request.routeId] = STATIC_LIMIT
			users[request.userId].requestLeft[request.routeId] = STATIC_LIMIT
			continue
		}
		users[request.userId].requestLeft[request.routeId] = min(limit, users[request.userId].requestLeft[request.routeId]+1)
	}
	fmt.Println("updated the limits")
}

func generateRandomNumber() int {
	return rand.Intn(1000 * 1000 * 1000)
}

type User struct {
	id          int
	limits      map[int]int
	requestLeft map[int]int
	mu          *sync.Mutex
}

func main() {
	rand.Seed(time.Now().UnixNano())

	algo := &SlidingAlgo{}
	algo.Initalize()
	rl := NewRateLimiter(algo)

	user := &User{
		
	}
	user = rl.AddAUser(*user)
	route := Route{
		name : "payment_create",
		endpoint: "payment/create",
		method: "POST",
	}
	
	route = rl.AddARoute(route)

	rl.UserRequest(user.id, route.id )
	rl.UserRequest(user.id, route.id )
	rl.UserRequest(user.id, route.id )
	rl.UserRequest(user.id, route.id )
	rl.UserRequest(user.id, route.id )
	rl.UserRequest(user.id, route.id )
	rl.UserRequest(user.id, route.id )
	rl.UserRequest(user.id, route.id )

	time.Sleep(2*time.Second)

	rl.UserRequest(user.id, route.id )
	rl.UserRequest(user.id, route.id )
	rl.UserRequest(user.id, route.id )
	select {}
}
