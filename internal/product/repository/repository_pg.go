package productrepo

import (
	"context"
	"fmt"

	"github.com/arfan21/vocagame/internal/entity"
	dbpostgres "github.com/arfan21/vocagame/pkg/db/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db    dbpostgres.Queryer
	rawDb *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db:    db,
		rawDb: db,
	}
}

func (r Repository) Begin(ctx context.Context) (tx pgx.Tx, err error) {
	return r.rawDb.Begin(ctx)
}

func (r Repository) WithTx(tx pgx.Tx) *Repository {
	r.db = tx
	return &r
}

func (r Repository) Create(ctx context.Context, data entity.Product) (err error) {
	query := `
		INSERT INTO products (user_id, name, description, stok, price)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = r.db.Exec(ctx, query,
		data.UserID,
		data.Name,
		data.Description,
		data.Stok,
		data.Price,
	)

	if err != nil {
		err = fmt.Errorf("product.repository.Create: failed to create product: %w", err)
		return
	}

	return
}
