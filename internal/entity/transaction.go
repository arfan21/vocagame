package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionStatus string

const (
	TransactionStatusProcessing TransactionStatus = "PROCESSING"
	TransactionStatusCompleted  TransactionStatus = "COMPLETED"
	TransactionStatusFailed     TransactionStatus = "FAILED"
)

type Transaction struct {
	ID                uuid.UUID           `json:"id"`
	UserID            uuid.UUID           `json:"user_id"`
	TransactionTypeID int                 `json:"transaction_type_id"`
	Status            TransactionStatus   `json:"status"`
	TotalAmount       decimal.Decimal     `json:"total_amount"`
	CreatedAt         string              `json:"created_at"`
	UpdatedAt         string              `json:"updated_at"`
	User              User                `json:"user"`
	TransactionType   TransactionType     `json:"transaction_type"`
	TransactionDetail []TransactionDetail `json:"transaction_detail"`
}

func (Transaction) TableName() string {
	return "transactions"
}

type TransactionType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (TransactionType) TableName() string {
	return "transaction_types"
}

type TransactionDetail struct {
	ID            uuid.UUID `json:"id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	ProductID     uuid.UUID `json:"product_id"`
	Qty           int       `json:"qty"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

func (TransactionDetail) TableName() string {
	return "transaction_details"
}
