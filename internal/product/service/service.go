package productsvc

import (
	"context"
	"fmt"

	"github.com/arfan21/vocagame/internal/entity"
	"github.com/arfan21/vocagame/internal/model"
	"github.com/arfan21/vocagame/internal/product"
	"github.com/arfan21/vocagame/pkg/validation"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	repo product.Repository
}

func New(repo product.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) WithTx(tx pgx.Tx) product.Service {
	s.repo = s.repo.WithTx(tx)
	return &s
}

func (s Service) Create(ctx context.Context, req model.ProductCreateRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("product.service.Create: failed to validate request : %w", err)
		return
	}

	data := entity.Product{
		UserID:      req.UserID,
		Name:        req.Name,
		Description: req.Description,
		Stok:        req.Stok,
		Price:       req.Price,
	}

	err = s.repo.Create(ctx, data)
	if err != nil {
		err = fmt.Errorf("product.service.Create: failed to create new product : %w", err)
		return
	}

	return nil
}
