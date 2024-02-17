package walletsvc

import (
	"context"
	"fmt"

	"github.com/arfan21/vocagame/internal/entity"
	"github.com/arfan21/vocagame/internal/model"
	"github.com/arfan21/vocagame/internal/wallet"
	"github.com/arfan21/vocagame/pkg/validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	repo wallet.Repository
}

func New(repo wallet.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) WithTx(tx pgx.Tx) wallet.Service {
	s.repo = s.repo.WithTx(tx)
	return &s
}

func (s Service) Create(ctx context.Context, req model.CreateWalletRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("wallet.service.Create: failed to validate request : %w", err)
		return
	}

	data := entity.Wallet{
		UserID: req.UserID,
	}

	err = s.repo.Create(ctx, data)
	if err != nil {
		err = fmt.Errorf("wallet.service.Create: failed to create new wallet : %w", err)
		return
	}

	return nil
}

func (s Service) GetByUserID(ctx context.Context, userID uuid.UUID, isForUpdate bool) (data entity.Wallet, err error) {
	return s.repo.GetByUserID(ctx, userID, isForUpdate)
}

func (s Service) UpdateBalance(ctx context.Context, data entity.Wallet) (err error) {
	return s.repo.UpdateBalance(ctx, data)
}
