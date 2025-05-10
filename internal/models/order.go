package models

import (
	"github.com/google/uuid"
	"time"
)

type OrderType string

const (
	OrderTypeStar    OrderType = "star"
	OrderTypePremium OrderType = "premium"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusCompleted OrderStatus = "completed"
	StatusFailed    OrderStatus = "failed"
)

type Order struct {
	ID            uuid.UUID   `json:"id" db:"id"`
	Type          OrderType   `json:"type" db:"type"`
	Status        OrderStatus `json:"status" db:"status"`
	Username      string      `json:"username" db:"username"`
	RecipientHash string      `json:"recipient_hash"`
	Quantity      *int        `json:"quantity" db:"quantity"`
	Months        *int        `json:"months,omitempty"`
	Amount        float64     `json:"amount" db:"amount"`
	WalletType    string      `json:"wallet_type" db:"wallet_type"`
	TxHash        *string     `json:"tx_hash" db:"tx_hash"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	CompletedAt   *time.Time  `json:"completed_at" db:"completed_at"`
	ErrorMessage  string      `json:"error_message" db:"error_message"`
}
