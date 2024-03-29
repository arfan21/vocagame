package product

import (
	"context"

	"github.com/arfan21/vocagame/internal/entity"
	productrepo "github.com/arfan21/vocagame/internal/product/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Begin(ctx context.Context) (tx pgx.Tx, err error)
	WithTx(tx pgx.Tx) *productrepo.Repository

	Create(ctx context.Context, data entity.Product) (err error)
	GetProducts(ctx context.Context, filter entity.ListProductFilter) (result []entity.Product, err error)
	GetTotalProduct(ctx context.Context, filter entity.ListProductFilter) (result int, err error)
	Update(ctx context.Context, data entity.Product) (err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	ReduceStok(ctx context.Context, id uuid.UUID, reduceBy int) (err error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) (result map[uuid.UUID]entity.Product, err error)
}
