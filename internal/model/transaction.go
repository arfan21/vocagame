package model

import (
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
