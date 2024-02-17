package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v4"
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
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	User              User                `json:"user"`
	TransactionType   TransactionType     `json:"transaction_type"`
	TransactionDetail []TransactionDetail `json:"transaction_detail"`
}

func (Transaction) TableName() string {
	return "transactions"
}

type TransactionType struct {
	ID   null.Int    `json:"id"`
	Name null.String `json:"name"`
}

func (TransactionType) TableName() string {
	return "transaction_types"
}

type TransactionDetail struct {
	ID            uuid.NullUUID      `json:"id"`
	TransactionID uuid.NullUUID      `json:"transaction_id"`
	ProductID     uuid.NullUUID      `json:"product_id"`
	Qty           null.Int           `json:"qty"`
	CreatedAt     null.Time          `json:"created_at"`
	UpdatedAt     null.Time          `json:"updated_at"`
	Product       TransactionProduct `json:"product"`
}

func (TransactionDetail) TableName() string {
	return "transaction_details"
}

type TransactionProduct struct {
	Name  null.String         `json:"name"`
	Price decimal.NullDecimal `json:"price"`
}
