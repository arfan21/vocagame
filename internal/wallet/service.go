package wallet

import (
	"context"

	"github.com/arfan21/vocagame/internal/model"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	WithTx(tx pgx.Tx) Service

	Create(ctx context.Context, req model.CreateWalletRequest) (err error)
}
