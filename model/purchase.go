package model

import "time"

// PurchaseData represents the data structure for a purchase
type PurchaseData struct {
	Data          time.Time `json:"data"`
	User          string    `json:"cliente"`
	UserData      string    `json:"dados"`
	ProductName   string    `json:"nome_produto"`
	Price         string    `json:"valor"`
	PaymentMethod string    `json:"pagamento"`
	UserInfo      string    `json:"obs"`
}
