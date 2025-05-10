package models

type StarOrderResponse struct {
	OrderID     string  `json:"order_id"`
	Status      string  `json:"status"`
	Username    string  `json:"username"`
	Quantity    int     `json:"quantity"`
	Amount      float64 `json:"amount"`
	CreatedAt   string  `json:"created_at"`
	CompletedAt *string `json:"completed_at,omitempty"`
	TxHash      *string `json:"tx_hash,omitempty"`
}

type PremiumOrderResponse struct {
	OrderID     string  `json:"order_id"`
	Status      string  `json:"status"`
	Username    string  `json:"username"`
	Months      int     `json:"months"`
	Amount      float64 `json:"amount"`
	CreatedAt   string  `json:"created_at"`
	CompletedAt *string `json:"completed_at,omitempty"`
	TxHash      *string `json:"tx_hash,omitempty"`
}
