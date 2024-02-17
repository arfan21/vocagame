package wallet

import (
	"context"

	"github.com/arfan21/vocagame/internal/entity"
	walletrepo "github.com/arfan21/vocagame/internal/wallet/repository"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Begin(ctx context.Context) (tx pgx.Tx, err error)
	WithTx(tx pgx.Tx) *walletrepo.Repository

	Create(ctx context.Context, data entity.Wallet) (err error)
}
