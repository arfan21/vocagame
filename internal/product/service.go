package product

import (
	"context"

	"github.com/arfan21/vocagame/internal/model"
	"github.com/arfan21/vocagame/pkg/pkgutil"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	WithTx(tx pgx.Tx) Service

	Create(ctx context.Context, req model.ProductCreateRequest) (err error)
	GetProducts(ctx context.Context, req model.GetListProductRequest) (res pkgutil.PaginationResponse, err error)
}
