package wallet

import (
	"context"

	"github.com/arfan21/vocagame/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	WithTx(tx pgx.Tx) Service

	Create(ctx context.Context, req model.CreateWalletRequest) (err error)
	GetByUserID(ctx context.Context, userID uuid.UUID, isForUpdate bool) (res model.WalletResponse, err error)
	UpdateBalance(ctx context.Context, req model.UpdateBalanceRequest) (err error)
}
