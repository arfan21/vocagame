package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID          uuid.UUID       `json:"id"`
	UserID      uuid.UUID       `json:"user_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Stok        int             `json:"stok"`
	Price       decimal.Decimal `json:"price"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	User        User            `json:"user"`
}

func (Product) TableName() string {
	return "products"
}

type ListProductFilter struct {
	UserID        uuid.NullUUID `jsonL:"user_id"`
	Name          string        `query:"name" json:"name"`
	Page          int           `query:"page" json:"page" validate:"min=1"`
	Limit         int           `query:"limit" json:"limit" validate:"min=1"`
	DisableOffset bool          `json:"-"`
}
