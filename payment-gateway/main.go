package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// --- ENUMS & CONSTANTS ---

type PaymentStatus string
type RefundStatus string

const (
	PaymentInitiated PaymentStatus = "INITIATED"
	PaymentSuccess   PaymentStatus = "SUCCESS"
	PaymentFailed    PaymentStatus = "FAILED"

	RefundInitiated RefundStatus = "INITIATED"
	RefundSuccess   RefundStatus = "SUCCESS"
	RefundFailed    RefundStatus = "FAILED"
)

// --- INTERFACES ---

// PaymentMethod abstracts all payment types (UPI, Card, Wallet, etc.)
type PaymentMethod interface {
	Process(payment *Payment) error
	MethodName() string
}

// --- PAYMENT METHODS ---

// UPI payment method implementation
// Demonstrates polymorphism via interface

type UPI struct{}

func (u *UPI) Process(payment *Payment) error {
	fmt.Println("Processing UPI payment for PaymentID:", payment.ID)
	// Simulate processing delay
	time.Sleep(100 * time.Millisecond)
	return nil // Always success for demo
}
func (u *UPI) MethodName() string { return "UPI" }

// Card payment method implementation
type Card struct{}

func (c *Card) Process(payment *Payment) error {
	fmt.Println("Processing Card payment for PaymentID:", payment.ID)
	time.Sleep(150 * time.Millisecond)
	return nil
}
func (c *Card) MethodName() string { return "Card" }

// Wallet payment method implementation
type Wallet struct{}

func (w *Wallet) Process(payment *Payment) error {
	fmt.Println("Processing Wallet payment for PaymentID:", payment.ID)
	time.Sleep(120 * time.Millisecond)
	return nil
}
func (w *Wallet) MethodName() string { return "Wallet" }

// --- DOMAIN MODELS ---

// Merchant represents a merchant in the system
type Merchant struct {
	ID   string
	Name string
}

// Payment represents a payment attempt
// Demonstrates encapsulation

type Payment struct {
	ID             string
	MerchantID     string
	Amount         float64
	Method         PaymentMethod
	Status         PaymentStatus
	IdempotencyKey string
	CreatedAt      time.Time
}

// Refund represents a refund request
type Refund struct {
	ID        string
	PaymentID string
	Amount    float64
	Status    RefundStatus
	CreatedAt time.Time
}

// --- PAYMENT GATEWAY ---

// PaymentGateway is the main orchestrator
// Demonstrates composition, thread safety, and OOP design

type PaymentGateway struct {
	merchants       map[string]*Merchant
	payments        map[string]*Payment
	refunds         map[string]*Refund
	idempotencyKeys map[string]string

	merchantsMu     sync.RWMutex
	paymentsMu      sync.RWMutex
	refundsMu       sync.RWMutex
	idempotencyMu   sync.RWMutex

	// Add a condition variable per idempotency key
	idempotencyConds map[string]*sync.Cond
}

// NewPaymentGateway initializes the gateway
func NewPaymentGateway() *PaymentGateway {
	return &PaymentGateway{
		merchants:         make(map[string]*Merchant),
		payments:          make(map[string]*Payment),
		refunds:           make(map[string]*Refund),
		idempotencyKeys:   make(map[string]string),
		idempotencyConds:  make(map[string]*sync.Cond),
	}
}

// RegisterMerchant adds a new merchant to the system
func (pg *PaymentGateway) RegisterMerchant(id, name string) {
	pg.merchantsMu.Lock()
	defer pg.merchantsMu.Unlock()
	pg.merchants[id] = &Merchant{ID: id, Name: name}
}

// InitiatePayment creates and processes a payment
func (pg *PaymentGateway) InitiatePayment(merchantID string, amount float64, method PaymentMethod, idempotencyKey string) (*Payment, error) {
	// 1. Atomically check and set idempotency key
	pg.idempotencyMu.Lock()
	if pid, exists := pg.idempotencyKeys[idempotencyKey]; exists {
		// If the payment is not yet available, wait for it
		cond, ok := pg.idempotencyConds[idempotencyKey]
		if !ok {
			// Create a new cond if not present (should not happen, but safe)
			cond = sync.NewCond(&pg.idempotencyMu)
			pg.idempotencyConds[idempotencyKey] = cond
		}
		for pid == "" {
			cond.Wait()
			pid = pg.idempotencyKeys[idempotencyKey]
		}
		pg.idempotencyMu.Unlock()
		pg.paymentsMu.RLock()
		payment := pg.payments[pid]
		pg.paymentsMu.RUnlock()
		return payment, nil
	}
	// Reserve the idempotency key with a placeholder (to block others)
	pg.idempotencyKeys[idempotencyKey] = "" // placeholder
	// Create a new cond for this idempotency key
	cond := sync.NewCond(&pg.idempotencyMu)
	pg.idempotencyConds[idempotencyKey] = cond
	pg.idempotencyMu.Unlock()

	// 2. Check merchant
	pg.merchantsMu.RLock()
	merchant, ok := pg.merchants[merchantID]
	pg.merchantsMu.RUnlock()
	if !ok {
		// Clean up the idempotency key reservation and cond
		pg.idempotencyMu.Lock()
		delete(pg.idempotencyKeys, idempotencyKey)
		if c, ok := pg.idempotencyConds[idempotencyKey]; ok {
			c.Broadcast()
			delete(pg.idempotencyConds, idempotencyKey)
		}
		pg.idempotencyMu.Unlock()
		return nil, errors.New("merchant not found")
	}

	// 3. Create and process payment
	payment := &Payment{
		ID:             generateID(),
		MerchantID:     merchant.ID,
		Amount:         amount,
		Method:         method,
		Status:         PaymentInitiated,
		IdempotencyKey: idempotencyKey,
		CreatedAt:      time.Now(),
	}
	err := method.Process(payment)
	if err != nil {
		payment.Status = PaymentFailed
	} else {
		payment.Status = PaymentSuccess
	}

	// 4. Store payment and update idempotency key
	pg.paymentsMu.Lock()
	pg.payments[payment.ID] = payment
	pg.paymentsMu.Unlock()

	pg.idempotencyMu.Lock()
	pg.idempotencyKeys[idempotencyKey] = payment.ID
	// Broadcast to all waiting goroutines that the payment is ready
	if c, ok := pg.idempotencyConds[idempotencyKey]; ok {
		c.Broadcast()
		delete(pg.idempotencyConds, idempotencyKey)
	}
	pg.idempotencyMu.Unlock()

	return payment, nil
}

// RequestRefund creates a refund for a successful payment
func (pg *PaymentGateway) RequestRefund(paymentID string, amount float64) (*Refund, error) {
	pg.paymentsMu.RLock()
	payment, ok := pg.payments[paymentID]
	pg.paymentsMu.RUnlock()
	if !ok || payment.Status != PaymentSuccess {
		return nil, errors.New("invalid or non-successful payment")
	}

	refund := &Refund{
		ID:        generateID(),
		PaymentID: paymentID,
		Amount:    amount,
		Status:    RefundInitiated,
		CreatedAt: time.Now(),
	}
	refund.Status = RefundSuccess

	pg.refundsMu.Lock()
	pg.refunds[refund.ID] = refund
	pg.refundsMu.Unlock()

	return refund, nil
}

// GetPaymentStatus fetches the status of a payment by ID
func (pg *PaymentGateway) GetPaymentStatus(paymentID string) (PaymentStatus, error) {
	pg.paymentsMu.RLock()
	defer pg.paymentsMu.RUnlock()
	if payment, ok := pg.payments[paymentID]; ok {
		return payment.Status, nil
	}
	return "", errors.New("payment not found")
}

// GetRefundStatus fetches the status of a refund by ID
func (pg *PaymentGateway) GetRefundStatus(refundID string) (RefundStatus, error) {
	pg.refundsMu.RLock()
	defer pg.refundsMu.RUnlock()
	if refund, ok := pg.refunds[refundID]; ok {
		return refund.Status, nil
	}
	return "", errors.New("refund not found")
}

// --- UTILS ---

var idCounter = 0
var idMu sync.Mutex

// generateID returns a unique string ID (for demo purposes)
func generateID() string {
	idMu.Lock()
	defer idMu.Unlock()
	idCounter++
	return fmt.Sprintf("id-%d", idCounter)
}

// --- MAIN (Example Usage) ---

func main() {
	gateway := NewPaymentGateway()
	gateway.RegisterMerchant("m1", "Merchant One")
	gateway.RegisterMerchant("m2", "Merchant Two")

	// Initiate a UPI payment
	payment1, err := gateway.InitiatePayment("m1", 100.0, &UPI{}, "idem-key-123")
	if err != nil {
		fmt.Println("Payment error:", err)
		return
	}
	fmt.Printf("Payment1: %+v\n", payment1)

	// Try duplicate payment with same idempotency key (should not create new payment)
	payment1Dup, _ := gateway.InitiatePayment("m1", 100.0, &UPI{}, "idem-key-123")
	fmt.Printf("Duplicate Payment1 (idempotency): %+v\n", payment1Dup)

	// Initiate a Card payment
	payment2, err := gateway.InitiatePayment("m2", 250.0, &Card{}, "idem-key-456")
	if err != nil {
		fmt.Println("Payment error:", err)
		return
	}
	fmt.Printf("Payment2: %+v\n", payment2)

	// Request a refund for payment1
	refund, err := gateway.RequestRefund(payment1.ID, 100.0)
	if err != nil {
		fmt.Println("Refund error:", err)
		return
	}
	fmt.Printf("Refund: %+v\n", refund)

	// Fetch payment status
	status, _ := gateway.GetPaymentStatus(payment1.ID)
	fmt.Println("Payment1 Status:", status)

	// Fetch refund status
	rstatus, _ := gateway.GetRefundStatus(refund.ID)
	fmt.Println("Refund Status:", rstatus)

	// Try refunding a failed payment (should error)
	_, err = gateway.RequestRefund("non-existent-id", 50.0)
	if err != nil {
		fmt.Println("Expected error for invalid refund:", err)
	}
} 