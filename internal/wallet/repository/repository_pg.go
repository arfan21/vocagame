package walletrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/arfan21/vocagame/internal/entity"
	"github.com/arfan21/vocagame/pkg/constant"
	dbpostgres "github.com/arfan21/vocagame/pkg/db/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db    dbpostgres.Queryer
	rawDb dbpostgres.Raw
}

func New(raw dbpostgres.Raw, queryer dbpostgres.Queryer) *Repository {
	return &Repository{
		db:    queryer,
		rawDb: raw,
	}
}

func (r Repository) Begin(ctx context.Context) (tx pgx.Tx, err error) {
	return r.rawDb.Begin(ctx)
}

func (r Repository) WithTx(tx pgx.Tx) *Repository {
	r.db = tx
	return &r
}

func (r Repository) Create(ctx context.Context, data entity.Wallet) (err error) {
	query := `
		INSERT INTO wallets (balance, user_id)
		VALUES ($1, $2)
	`

	_, err = r.db.Exec(ctx, query,
		data.Balance,
		data.UserID,
	)

	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLUniqueViolation {
				err = constant.ErrWalletAlreadyCreated
			}
		}

		err = fmt.Errorf("wallet.repository.Create: failed to create wallet: %w", err)
		return
	}

	return
}

func (r Repository) GetByUserID(ctx context.Context, userID uuid.UUID, isForUpdate bool) (data entity.Wallet, err error) {
	query := `
		SELECT id, user_id, balance, created_at, updated_at
		FROM wallets
		WHERE user_id = $1
	`

	if isForUpdate {
		query += " FOR UPDATE"
	}

	err = r.db.QueryRow(ctx, query, userID).Scan(
		&data.ID,
		&data.UserID,
		&data.Balance,
		&data.CreatedAt,
		&data.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrWalletNotFound
		} else {
			err = fmt.Errorf("wallet.repository.GetByUserID: failed to get wallet: %w", err)
		}
		return
	}

	return
}

func (r Repository) UpdateBalance(ctx context.Context, data entity.Wallet) (err error) {
	query := `
		UPDATE wallets
		SET balance = $1
		WHERE id = $2
	`

	_, err = r.db.Exec(ctx, query, data.Balance, data.ID)
	if err != nil {
		err = fmt.Errorf("wallet.repository.UpdateBalance: failed to update balance: %w", err)
		return
	}

	return
}
