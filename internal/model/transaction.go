package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateDepositTransactionRequest struct {
	UserID uuid.UUID       `json:"user_id" validate:"required"`
	Amount decimal.Decimal `json:"amount" validate:"required,dgt=0"`
}

type CreateWithdrawTransactionRequest struct {
	UserID uuid.UUID       `json:"user_id" validate:"required"`
	Amount decimal.Decimal `json:"amount" validate:"required"`
}

type CreateTransactionResponse struct {
	TransactionID string `json:"transaction_id"`
}

type GetTransactionByIDRequest struct {
	ID     uuid.UUID `json:"id" validate:"required"`
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

type GetTransactionResponse struct {
	ID              uuid.UUID                   `json:"id"`
	UserID          uuid.UUID                   `json:"user_id"`
	TransactionType string                      `json:"transaction_type"`
	Status          string                      `json:"status"`
	TotalAmount     decimal.Decimal             `json:"total_amount"`
	CreatedAt       time.Time                   `json:"created_at"`
	UpdatedAt       time.Time                   `json:"updated_at"`
	Details         []TransactionDetailResponse `json:"details,omitempty"`
}

type TransactionDetailResponse struct {
	ID           uuid.UUID       `json:"id"`
	ProductID    uuid.UUID       `json:"product_id"`
	Qty          int             `json:"qty"`
	ProductName  string          `json:"product_name"`
	ProductPrice decimal.Decimal `json:"product_price"`
}

type CheckoutTransactionRequest struct {
	UserID   uuid.UUID                `json:"user_id" validate:"required"`
	Products []CheckoutProductRequest `json:"products" validate:"required,min=1,dive,required"`
}

type CheckoutProductRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Qty       int       `json:"qty" validate:"required,min=1"`
}
