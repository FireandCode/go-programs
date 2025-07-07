package main

import (
	"sync"
	"testing"
)

func TestConcurrentInitiatePayment_Idempotency(t *testing.T) {
	gateway := NewPaymentGateway()
	merchantID := "m1"
	gateway.RegisterMerchant(merchantID, "Merchant One")

	idempotencyKey := "concurrent-key-1"
	concurrency := 20

	// To collect all returned payment pointers
	results := make([]*Payment, concurrency)
	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(idx int) {
			defer wg.Done()
			payment, err := gateway.InitiatePayment(merchantID, 100.0, &UPI{}, idempotencyKey)
			if err != nil {
				t.Fatalf("goroutine %d: unexpected error: %v", idx, err)
			}
			results[idx] = payment
		}(i)
	}
	wg.Wait()

	// Ensure none of the results are nil
	for i, p := range results {
		if p == nil {
			t.Fatalf("payment at index %d is nil", i)
		}
	}

	// All returned payments should have the same ID
	firstID := results[0].ID
	for i, p := range results {
		if p.ID != firstID {
			t.Errorf("payment ID mismatch at index %d: got %s, want %s", i, p.ID, firstID)
		}
	}

	// Only one payment should exist in the payments map for this idempotency key
	gateway.idempotencyMu.RLock()
	paymentID := gateway.idempotencyKeys[idempotencyKey]
	gateway.idempotencyMu.RUnlock()

	gateway.paymentsMu.RLock()
	_, exists := gateway.payments[paymentID]
	gateway.paymentsMu.RUnlock()
	if !exists {
		t.Errorf("expected payment with ID %s to exist in payments map", paymentID)
	}
} 