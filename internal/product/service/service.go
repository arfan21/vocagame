package productsvc

import (
	"context"
	"fmt"

	"github.com/arfan21/vocagame/internal/entity"
	"github.com/arfan21/vocagame/internal/model"
	"github.com/arfan21/vocagame/internal/product"
	"github.com/arfan21/vocagame/pkg/constant"
	"github.com/arfan21/vocagame/pkg/pkgutil"
	"github.com/arfan21/vocagame/pkg/validation"
	"github.com/google/uuid"
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

func (s Service) getProducts(ctx context.Context, filter entity.ListProductFilter) (res pkgutil.PaginationResponse[[]model.GetProductResponse], err error) {
	results, err := s.repo.GetProducts(ctx, filter)
	if err != nil {
		err = fmt.Errorf("product.service.GetProducts: failed to get products from db : %w", err)
		return
	}

	var resData []model.GetProductResponse

	resData = make([]model.GetProductResponse, len(results))

	for i, result := range results {
		resData[i].ID = result.ID
		resData[i].Name = result.Name
		resData[i].Description = result.Description
		resData[i].Stok = result.Stok
		resData[i].Price = result.Price
		resData[i].OwnerID = result.User.ID
		resData[i].OwnerName = result.User.Fullname
	}

	total, err := s.repo.GetTotalProduct(ctx, filter)
	if err != nil {
		err = fmt.Errorf("product.service.GetProducts: failed to get total product from db : %w", err)
		return
	}

	totalPage := 0
	if total%filter.Limit != 0 {
		totalPage = total/filter.Limit + 1
	} else {
		totalPage = total / filter.Limit
	}

	res = pkgutil.PaginationResponse[[]model.GetProductResponse]{
		TotalData: total,
		TotalPage: totalPage,
		Page:      filter.Page,
		Limit:     filter.Limit,
		Data:      resData,
	}

	return
}

func (s Service) GetProducts(ctx context.Context, req model.GetListProductRequest) (res pkgutil.PaginationResponse[[]model.GetProductResponse], err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("product.service.GetProducts: failed to validate request : %w", err)
		return
	}

	filter := entity.ListProductFilter{
		Name:   req.Name,
		Page:   req.Page,
		Limit:  req.Limit,
		UserID: req.OwnerID,
		ID:     req.ProductID,
	}

	return s.getProducts(ctx, filter)
}

func (s Service) Update(ctx context.Context, req model.ProductUpdateRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("product.service.Update: failed to validate request : %w", err)
		return
	}

	// check if product exist
	resultProduct, err := s.getProducts(ctx, entity.ListProductFilter{
		ID:    uuid.NullUUID{UUID: req.ID, Valid: true},
		Limit: 1,
		Page:  1,
	})
	if err != nil {
		err = fmt.Errorf("product.service.Update: failed to get product : %w", err)
		return
	}

	if len(resultProduct.Data) == 0 {
		err = constant.ErrProductNotFound
		return
	}

	// check if user is owner of product
	if resultProduct.Data[0].OwnerID != req.UserID {
		err = constant.ErrCannotUpdateNotOwner
		return
	}

	data := entity.Product{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Stok:        req.Stok,
		Price:       req.Price,
	}

	err = s.repo.Update(ctx, data)
	if err != nil {
		err = fmt.Errorf("product.service.Update: failed to update product : %w", err)
		return
	}

	return
}

func (s Service) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) (err error) {
	// check if product exist
	resultProduct, err := s.getProducts(ctx, entity.ListProductFilter{
		ID:    uuid.NullUUID{UUID: id, Valid: true},
		Limit: 1,
		Page:  1,
	})
	if err != nil {
		err = fmt.Errorf("product.service.Delete: failed to get product : %w", err)
		return
	}

	if len(resultProduct.Data) == 0 {
		err = constant.ErrProductNotFound
		return
	}

	// check if user is owner of product
	if resultProduct.Data[0].OwnerID != userID {
		err = constant.ErrCannotDeleteNotOwner
		return
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		err = fmt.Errorf("product.service.Delete: failed to delete product : %w", err)
		return
	}

	return
}
