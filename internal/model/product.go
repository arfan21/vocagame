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

type GetListProductRequest struct {
	Name      string        `query:"name" json:"name"`
	Page      int           `query:"page" json:"page" validate:"min=1"`
	Limit     int           `query:"limit" json:"limit" validate:"min=1"`
	OwnerID   uuid.NullUUID `query:"owner_id" json:"owner_id"`
	ProductID uuid.NullUUID `query:"product_id" json:"product_id"`
}

type GetProductResponse struct {
	ID        uuid.UUID       `json:"id" swaggertype:"string"`
	Name      string          `json:"name"`
	Stok      int             `json:"stok"`
	Price     decimal.Decimal `json:"price" swaggertype:"string"`
	OwnerID   uuid.UUID       `json:"owner_id" swaggertype:"string"`
	OwnerName string          `json:"owner_name"`
}
