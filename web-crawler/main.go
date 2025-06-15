package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

/*
	searchEngine struct {
 Depth *int
 rateLimiter RateLimiter
 webcrawlers sync.Map WebCrawler
 details map[request_id]sync.Map // key as the requestId and value as the details
 mut *sync.Mutex
}

WebCrawler struct {
	id int
	requestId int
	currentDepth int
}


type Request {
	id int
	isThrottled bool
}
type User {
	id int
	requests sync.Map
}


RateLimiter struct -
	requestId int
	maxLimit
	limits sync.Map
	internal 5 sec

	  Details struct -
	- requestId int
	- seedURL string
	- page titles sync.Map
	- meta descriptions sync.Map
	- URLlinks sync.Map //handle duplicate



without wait groups. without channels

concurrency
	- go routines
	- why channel will be required and where can i use them here?
	  - i am using sync.Map for the shared resources.

*/

type SearchEngine struct {
	depth int 
	requests *sync.Map
	users map[int]User
	RateLimiter RateLimiter
	details *sync.Map
	requestch chan Request
}

type Request struct {
	id int
	url string 
	userId int 
	isThrottled bool
}

type User struct {
	id int
	name string 
}


type RateLimiter struct {
	maxLimit int 
	mut map[int]*sync.Mutex
	limits *sync.Map // user_id -> limits 
	interval time.Duration
}

type Website struct {
	url string 
	title string 
	description string 
}
type Details struct {
	requestId int 
	webSiteDetails *sync.Map 
}

/*
rateLImiter struct
check()

user struct
request struct
details struct

webCrawler struct
crawl
fetchURL
fetchTitleAndDescription
FetchLinks()

searchEngine struct
NewSearchEngine()
AddUser()
AddARequest()
SearchAURL()
ShowDetails()
ShowDetailsForARequestId()
updateRateLimits


type SearchEngine struct {
	depth int 
	requests sync.Map
	users sync.Map
	RateLimiter RateLimiter
	details map[int]Details
}
*/

func (s *SearchEngine) updateRateLimits(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return // Exit when context is canceled
		case <-ticker.C:
			s.RateLimiter.limits.Range(func(key, value any) bool {
				userId := key.(int)
				currentLimit := value.(int)
				s.RateLimiter.limits.Store(userId, max(s.RateLimiter.maxLimit, currentLimit+1))
				return true
			})
		}
	}
}

func(s *SearchEngine) ShowDetailsForARequestId(requestId int ) {
details , ok := s.details.Load(requestId)

if !ok {
	log.Println("details are not for this request")
	return 
}
detailsStruct := details.(Details)
detailsStruct.webSiteDetails.Range(func (key ,value any ) bool {
			website := value.(Website)
			log.Printf("url %s, title %s, description %s \n", website.url, website.title, website.description)
			return true
		})
}

func(s *SearchEngine) ShowDetails() {

	s.details.Range(func (key, value any ) bool {
		detailsStruct := value.(Details)
		log.Println("RequestId ", detailsStruct.requestId)
		
		detailsStruct.webSiteDetails.Range(func (k ,v any ) bool {
			website := v.(Website)
			log.Printf("url %s, title %s, description %s \n", website.url, website.title, website.description)
			return true
		})
		return true 
	})

}

func(s *SearchEngine) SearchAURL(request Request) (bool, error) {
	s.requests.Store(request.id, request)
	pass := s.RateLimiter.Check(request.userId)
	if !pass {
		request.isThrottled = true 
		s.requests.Store(request.id, request)
		return false, errors.New("Request is throttled, please try again after sometime")
	} 
	details := Details{
		requestId: request.id,
		webSiteDetails: &sync.Map{},
	}
	
	s.details.Store(request.id, details)
	log.Println("searching the URL")
	Crawl(request.url, 0, request.id, s)
	log.Println("search successful below are the results")
	detail, _ := s.details.Load(request.id)

	detail.(Details).webSiteDetails.Range(func (key ,value any ) bool {
		website := value.(Website)
		log.Printf("url %s, title %s, description %s \n", website.url, website.title, website.description)
		return true
	})
	
	return true, nil
}

var randGen = rand.New(rand.NewSource(time.Now().UnixMilli()))
func generateRandomInt() int {
	digits := 8
    if digits <= 0 {
        return 0 // Return 0 for invalid input
    }

    min := int(math.Pow10(digits - 1))
    max := int(math.Pow10(digits)) - 1
    return min + randGen.Intn(max-min+1)
}

func CreateARequest(url string, userId int ) Request {
	id := generateRandomInt()
	request := Request{
		id: id,
		url: url,
		userId:  userId,
		isThrottled: false,
	}
	return request
}

func(s *SearchEngine) ProcessRequest() {
	log.Println("processing request")
		for request := range  s.requestch{
			log.Println("Received request:", request)
			s.SearchAURL(request)
		}
		wg.Done()
}

func(s *SearchEngine) AddAUser(name string) User{
	id := generateRandomInt()
	user := User{
		id: id,
		name: name,
	}
	s.users[id] = user

	return user
}

func NewRateLimiter(maxLimit int, interval time.Duration) RateLimiter {
	limits := &sync.Map{}

	return RateLimiter{
		maxLimit: maxLimit,
		limits: limits,
		interval: interval,
	}
}
func NewSearchEngine(depth int, rl RateLimiter, ctx context.Context) *SearchEngine  {
	sEng := SearchEngine{}
	sEng.depth = depth
	sEng.requests = &sync.Map{}
	sEng.users = make(map[int]User)
	sEng.requestch = make(chan Request, 10)

	sEng.RateLimiter = rl 
	sEng.details = &sync.Map{}

	go sEng.updateRateLimits(ctx)

	return &sEng
}

func Crawl(url string, depth int,  requestId int,  s *SearchEngine)  {
	detail, _ := s.details.Load(requestId)
	
	_, ok := detail.(Details).webSiteDetails.Load(url)

	if ok {
		return ;
	}
	website := Website{
		url: url,
	}
	wg.Add(1)
	defer wg.Done()

	detail.(Details).webSiteDetails.LoadOrStore(url, website)

	details := fetchURLDetails(url)
	title, description := fetchTitleAndDescription(details)

	website.title = title
	website.description = description
	detail.(Details).webSiteDetails.LoadOrStore(url, website)

	urls := fetchLinks(url)

	for _, url := range urls {
		if _, ok := detail.(Details).webSiteDetails.Load(url); ok || depth+1 > s.depth {
			continue
		}
		wg.Add(1)
		go func(nextURL string ) {
			defer wg.Done()
			Crawl(nextURL, depth+1, requestId, s)
		}(url)
	}
}

func  fetchURLDetails(url string) string {
		return "URL Details"
}

func fetchTitleAndDescription(details string) (string, string) {
	return "Hello", "World"
}

func  fetchLinks(url string) ([]string) {
	urls := []string{"https://google.com","https://google.com", "https://example.com", "https://geeksforgeeks.com" }

	return urls
}

func(r *RateLimiter) Check(userId int) bool {
	r.mut[userId].Lock()
	defer r.mut[userId].Unlock()

	value, ok := r.limits.Load(userId)
	if !ok {
		value = r.maxLimit
	}
	if value == 0 {
		return false
	}

	valueInt, ok  := value.(int)
	if !ok {
		log.Println("Limit is not a int for userid: ", userId)
		return false
	}
	//race condition how do i handle it efficiently
	r.limits.LoadOrStore(userId, valueInt-1)

	return true 
}

var wg sync.WaitGroup

func main()  {
	rl := NewRateLimiter(4, 4*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	seng := NewSearchEngine(4,rl, ctx)
	user := seng.AddAUser("pradeep")
	request := CreateARequest("https://example.com", user.id)
	wg.Add(1)
	go seng.ProcessRequest()

	go func () {
	
		seng.requestch <- request
		seng.requestch <- request
		seng.requestch <- request
		seng.requestch <- request
		close(seng.requestch)
	
	}()

	fmt.Println(uuid.New())

	wg.Wait()
	cancel()
}

