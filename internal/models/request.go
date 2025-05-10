package models

type CreateStarOrderRequest struct {
	Username      string `json:"username" binding:"required"`
	RecipientHash string `json:"recipient_hash" binding:"required"`
	Quantity      int    `json:"quantity" binding:"required,min=50,max=1000000"`
	WalletType    string `json:"wallet_type" binding:"required"`
}

type CreatePremiumOrderRequest struct {
	Username      string `json:"username" binding:"required"`
	RecipientHash string `json:"recipient_hash" binding:"required"`
	Months        int    `json:"months" binding:"required,oneof=3 6 12"`
	WalletType    string `json:"wallet_type" binding:"required"`
}
