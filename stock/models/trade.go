package models

import (
	"errors"
	"time"
)

type TradeType string

const (
	Buy  TradeType = "buy"
	Sell TradeType = "sell"
)

type Trade struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Type      TradeType `gorm:"not null" json:"type"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Symbol    string    `gorm:"not null" json:"symbol"`
	Shares    int       `gorm:"not null" json:"shares"`
	Price     float64   `gorm:"not null" json:"price"`
	Timestamp int64     `gorm:"not null" json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Trade) Validate() error {
	if t.Type != Buy && t.Type != Sell {
		return errors.New("trade type must be either 'buy' or 'sell'")
	}

	if t.Shares < 1 || t.Shares > 100 {
		return errors.New("shares must be between 1 and 100")
	}

	if t.Price <= 0 {
		return errors.New("price must be greater than 0")
	}

	if t.Symbol == "" {
		return errors.New("symbol is required")
	}

	if t.UserID == 0 {
		return errors.New("user_id is required")
	}

	if t.Timestamp == 0 {
		t.Timestamp = time.Now().UnixMilli()
	}

	return nil
} 