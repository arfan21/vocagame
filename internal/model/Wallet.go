package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateWalletRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

type WalletResponse struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type UpdateBalanceRequest struct {
	ID      uuid.UUID       `json:"id" validate:"required"`
	Balance decimal.Decimal `json:"balance" validate:"required,dgt=0"`
	UserID  uuid.UUID       `json:"user_id" validate:"required"`
}
