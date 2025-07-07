package main

import "fmt"

type Observer interface {
	ID() string
	Update(msg string)
}

type NotificationService struct {
	observers map[string]Observer
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		observers: make(map[string]Observer),
	}
}

func (n *NotificationService) Register(o Observer) {
	n.observers[o.ID()] = o
}

func (n *NotificationService) Unregister(id string) {
	delete(n.observers, id)
}

func (n *NotificationService) NotifyAll(msg string) {
	for _, o := range n.observers {
		o.Update(msg)
	}
}


type NewsAgency struct {
	notificationService *NotificationService
}

func NewNewsAgency(n *NotificationService) *NewsAgency {
	return &NewsAgency{notificationService: n}
}

func (n *NewsAgency) Publish(news string) {
	fmt.Println("NewsAgency published:", news)
	n.notificationService.NotifyAll(news)
}


type Subscriber struct {
	id   string
	name string
}

func NewSubscriber(id, name string) *Subscriber {
	return &Subscriber{id: id, name: name}
}

func (s *Subscriber) ID() string {
	return s.id
}

func (s *Subscriber) Update(msg string) {
	fmt.Printf("[%s] received: %s\n", s.name, msg)
}

func main() {
	// Step 1: Create notification service
	notifications := NewNotificationService()

	// Step 2: Register observers
	alice := NewSubscriber("1", "Alice")
	bob := NewSubscriber("2", "Bob")
	notifications.Register(alice)
	notifications.Register(bob)

	// Step 3: Create NewsAgency with the notification service
	agency := NewNewsAgency(notifications)

	// Step 4: Publish news (automatically notifies all subscribers)
	agency.Publish("Go 1.22 Released!")
	agency.Publish("Observer Pattern Made Easy")

	// Step 5: Unsubscribe someone
	notifications.Unregister("2") // Bob unsubscribed
	agency.Publish("Only Alice will see this")
}
