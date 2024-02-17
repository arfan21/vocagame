package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID        uuid.UUID       `json:"id"`
	Balance   decimal.Decimal `json:"balance"`
	UserID    uuid.UUID       `json:"user_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	User      User            `json:"user"`
}

func (Wallet) TableName() string {
	return "wallets"
}
