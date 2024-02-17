package transactionsvc

import (
	"context"
	"fmt"

	"github.com/arfan21/vocagame/internal/entity"
	"github.com/arfan21/vocagame/internal/model"
	"github.com/arfan21/vocagame/internal/product"
	"github.com/arfan21/vocagame/internal/transaction"
	"github.com/arfan21/vocagame/internal/wallet"
	"github.com/arfan21/vocagame/pkg/constant"
	"github.com/arfan21/vocagame/pkg/validation"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Service struct {
	repo       transaction.Repository
	walletSvc  wallet.Service
	productSvc product.Service
}

func New(repo transaction.Repository, walletSvc wallet.Service, productSvc product.Service) *Service {
	return &Service{repo: repo, walletSvc: walletSvc, productSvc: productSvc}
}

func (s Service) CreateDepositTransaction(ctx context.Context, req model.CreateDepositTransactionRequest) (res model.CreateTransactionResponse, err error) {
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

	walletDataReq := model.UpdateBalanceRequest{
		ID:      walletData.ID,
		Balance: walletData.Balance,
		UserID:  walletData.UserID,
	}

	err = s.walletSvc.WithTx(tx).UpdateBalance(ctx, walletDataReq)
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

func (s Service) CreateWithdrawTransaction(ctx context.Context, req model.CreateWithdrawTransactionRequest) (res model.CreateTransactionResponse, err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("transaction.service.CreateWithdrawTransaction: failed to validate request: %w", err)
		return
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("transaction.service.CreateWithdrawTransaction: failed to begin transaction: %w", err)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}

		err = tx.Commit(ctx)
		if err != nil {
			err = fmt.Errorf("transaction.service.CreateWithdrawTransaction: failed to commit transaction: %w", err)
			return
		}
	}()

	walletData, err := s.walletSvc.WithTx(tx).GetByUserID(ctx, req.UserID, true)
	if err != nil {
		err = fmt.Errorf("transaction.service.CreateWithdrawTransaction: failed to get wallet data: %w", err)
		return
	}

	if walletData.Balance.LessThan(req.Amount) {
		err = constant.ErrInsufficientBalance
		return
	}

	walletData.Balance = walletData.Balance.Sub(req.Amount)

	walletDataReq := model.UpdateBalanceRequest{
		ID:      walletData.ID,
		Balance: walletData.Balance,
		UserID:  walletData.UserID,
	}

	err = s.walletSvc.WithTx(tx).UpdateBalance(ctx, walletDataReq)
	if err != nil {
		err = fmt.Errorf("transaction.service.CreateWithdrawTransaction: failed to update wallet balance: %w", err)
		return
	}

	transactionData := entity.Transaction{
		UserID:            req.UserID,
		TransactionTypeID: constant.TransactionTypeWithdrawID,
		Status:            entity.TransactionStatusCompleted,
		TotalAmount:       req.Amount,
	}

	idTx, err := s.repo.WithTx(tx).Create(ctx, transactionData)
	if err != nil {
		err = fmt.Errorf("transaction.service.CreateWithdrawTransaction: failed to create transaction: %w", err)
		return
	}

	res.TransactionID = idTx.String()

	return
}

func (s Service) GetHistoryWalletByUserID(ctx context.Context, userID uuid.UUID) (res []model.GetTransactionResponse, err error) {
	transactions, err := s.repo.GetHistoryWalletByUserID(ctx, userID)
	if err != nil {
		err = fmt.Errorf("transaction.service.GetHistoryWalletByUserID: failed to get history wallet: %w", err)
		return
	}

	res = make([]model.GetTransactionResponse, len(transactions))

	for i, transaction := range transactions {
		res[i].ID = transaction.ID
		res[i].UserID = transaction.UserID
		res[i].TransactionType = transaction.TransactionType.Name
		res[i].Status = string(transaction.Status)
		res[i].CreatedAt = transaction.CreatedAt
		res[i].UpdatedAt = transaction.UpdatedAt

		if transaction.TransactionTypeID == constant.TransactionTypeWithdrawID {
			res[i].TotalAmount = transaction.TotalAmount.Neg()
		} else {
			res[i].TotalAmount = transaction.TotalAmount
		}
	}

	return
}

func (s Service) Checkout(ctx context.Context, req model.CheckoutTransactionRequest) (res model.CreateTransactionResponse, err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("transaction.service.Checkout: failed to validate request: %w", err)
		return
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("transaction.service.Checkout: failed to begin transaction: %w", err)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}

		err = tx.Commit(ctx)
		if err != nil {
			err = fmt.Errorf("transaction.service.Checkout: failed to commit transaction: %w", err)
			return
		}
	}()

	productIds := make([]uuid.UUID, len(req.Products))

	for i, v := range req.Products {
		productIds[i] = v.ProductID
	}

	products, err := s.productSvc.WithTx(tx).GetByIDs(ctx, productIds)
	if err != nil {
		err = fmt.Errorf("transaction.service.Checkout: failed to get products: %w", err)
		return
	}

	if len(products) != len(req.Products) {
		productId := ""
		for _, v := range req.Products {
			if _, ok := products[v.ProductID]; !ok {
				productId = v.ProductID.String()
				break
			}
		}

		errProductNotFound := *constant.ErrProductNotFound
		errProductNotFound.Message = fmt.Sprintf("product with id '%s' not found", productId)
		err = &errProductNotFound
		return
	}

	productUpdateRequests := make([]model.ReduceStokRequest, len(req.Products))
	totalAmount := decimal.NewFromInt(0)

	// check stok
	for i, v := range req.Products {
		product := products[v.ProductID]
		if product.Stok < v.Qty {
			errProductStockNotEnough := *constant.ErrProductStokNotEnough
			errProductStockNotEnough.Message = fmt.Sprintf("product with name %s stok not enough", product.Name)
			err = &errProductStockNotEnough
			return
		}

		if product.OwnerID == req.UserID {
			err = constant.ErrCannotPurchaseOwnProduct
			return
		}

		productUpdateRequests[i] = model.ReduceStokRequest{
			ID:       product.ID,
			ReduceBy: v.Qty,
		}

		totalAmount = totalAmount.Add(product.Price.Mul(decimal.NewFromInt(int64(v.Qty))))
	}

	walletData, err := s.walletSvc.WithTx(tx).GetByUserID(ctx, req.UserID, true)
	if err != nil {
		err = fmt.Errorf("transaction.service.Checkout: failed to get wallet data: %w", err)
		return
	}

	if walletData.Balance.LessThan(totalAmount) {
		err = constant.ErrInsufficientBalance
		return
	}

	walletData.Balance = walletData.Balance.Sub(totalAmount)

	walletDataReq := model.UpdateBalanceRequest{
		ID:      walletData.ID,
		Balance: walletData.Balance,
		UserID:  walletData.UserID,
	}

	err = s.walletSvc.WithTx(tx).UpdateBalance(ctx, walletDataReq)
	if err != nil {
		err = fmt.Errorf("transaction.service.Checkout: failed to update wallet balance: %w", err)
		return
	}

	transactionData := entity.Transaction{
		UserID:            req.UserID,
		TransactionTypeID: constant.TransactionTypePurchaseID,
		Status:            entity.TransactionStatusCompleted,
		TotalAmount:       totalAmount,
	}

	idTx, err := s.repo.WithTx(tx).Create(ctx, transactionData)
	if err != nil {
		err = fmt.Errorf("transaction.service.Checkout: failed to create transaction: %w", err)
		return
	}

	transactionDetailData := make([]entity.TransactionDetail, len(req.Products))

	for i, v := range req.Products {
		transactionDetailData[i] = entity.TransactionDetail{
			TransactionID: idTx,
			ProductID:     v.ProductID,
			Qty:           v.Qty,
		}
	}

	err = s.repo.WithTx(tx).CreateDetail(ctx, transactionDetailData)
	if err != nil {
		err = fmt.Errorf("transaction.service.Checkout: failed to create transaction detail: %w", err)
		return
	}

	err = s.productSvc.WithTx(tx).BatchReduceStok(ctx, productUpdateRequests)
	if err != nil {
		err = fmt.Errorf("transaction.service.Checkout: failed to update product stok: %w", err)
		return
	}

	res.TransactionID = idTx.String()

	return
}
