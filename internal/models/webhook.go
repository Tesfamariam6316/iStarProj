package models

import "time"

type WebhookPayload struct {
	EventType   string                 `json:"event_type"`
	OccurredAt  time.Time              `json:"occurred_at"`
	Order       map[string]interface{} `json:"order"`
	TxHash      *string                `json:"tx_hash,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Quantity    *int                   `json:"quantity,omitempty"`
}
