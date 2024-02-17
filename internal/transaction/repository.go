package transaction

import (
	"context"

	"github.com/arfan21/vocagame/internal/entity"
	transactionrepo "github.com/arfan21/vocagame/internal/transaction/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Begin(ctx context.Context) (tx pgx.Tx, err error)
	WithTx(tx pgx.Tx) *transactionrepo.Repository

	Create(ctx context.Context, data entity.Transaction) (id uuid.UUID, err error)
	GetHistoryWalletByUserID(ctx context.Context, userID uuid.UUID) (res []entity.Transaction, err error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (res entity.Transaction, err error)
}
