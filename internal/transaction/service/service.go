package transactionsvc

import (
	"context"
	"fmt"

	"github.com/arfan21/vocagame/internal/entity"
	"github.com/arfan21/vocagame/internal/model"
	"github.com/arfan21/vocagame/internal/transaction"
	"github.com/arfan21/vocagame/internal/wallet"
	"github.com/arfan21/vocagame/pkg/constant"
	"github.com/arfan21/vocagame/pkg/validation"
)

type Service struct {
	repo      transaction.Repository
	walletSvc wallet.Service
}

func New(repo transaction.Repository, walletSvc wallet.Service) *Service {
	return &Service{repo: repo, walletSvc: walletSvc}
}

func (s *Service) CreateDepositTransaction(ctx context.Context, req model.CreateDepositTransactionRequest) (res model.CreateTransactionResponse, err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("transaction.service.CreateDepositTransaction: failed to validate request: %w", err)
		return
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("transaction.service.CreateDepositTransaction: failed to begin transaction: %w", err)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}

		err = tx.Commit(ctx)
		if err != nil {
			err = fmt.Errorf("transaction.service.CreateDepositTransaction: failed to commit transaction: %w", err)
			return
		}
	}()

	walletData, err := s.walletSvc.WithTx(tx).GetByUserID(ctx, req.UserID, true)
	if err != nil {
		err = fmt.Errorf("transaction.service.CreateDepositTransaction: failed to get wallet data: %w", err)
		return
	}

	walletData.Balance = walletData.Balance.Add(req.Amount)

	err = s.walletSvc.WithTx(tx).UpdateBalance(ctx, walletData)
	if err != nil {
		err = fmt.Errorf("transaction.service.CreateDepositTransaction: failed to update wallet balance: %w", err)
		return
	}

	transactionData := entity.Transaction{
		UserID:            req.UserID,
		TransactionTypeID: constant.TransactionTypeDepositID,
		Status:            entity.TransactionStatusCompleted,
		TotalAmount:       req.Amount,
	}

	idTx, err := s.repo.WithTx(tx).Create(ctx, transactionData)
	if err != nil {
		err = fmt.Errorf("transaction.service.CreateDepositTransaction: failed to create transaction: %w", err)
		return
	}

	res.TransactionID = idTx.String()

	return
}
