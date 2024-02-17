package model

import "github.com/google/uuid"

type CreateWalletRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}
