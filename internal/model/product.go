package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductCreateRequest struct {
	Name        string          `json:"name" validate:"required"`
	Stok        int             `json:"stok" validate:"required"`
	Description string          `json:"description" validate:"required"`
	Price       decimal.Decimal `json:"price" validate:"required" swaggertype:"string"`
	UserID      uuid.UUID       `json:"user_id" validate:"required"`
}
