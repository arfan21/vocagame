package transaction

import (
	"context"

	"github.com/arfan21/vocagame/internal/model"
	"github.com/google/uuid"
)

type Service interface {
	CreateDepositTransaction(ctx context.Context, req model.CreateDepositTransactionRequest) (res model.CreateTransactionResponse, err error)
	CreateWithdrawTransaction(ctx context.Context, req model.CreateWithdrawTransactionRequest) (res model.CreateTransactionResponse, err error)
	GetHistoryWalletByUserID(ctx context.Context, userID uuid.UUID) (res []model.GetTransactionResponse, err error)
	Checkout(ctx context.Context, req model.CheckoutTransactionRequest) (res model.CreateTransactionResponse, err error)
	GetByID(ctx context.Context, req model.GetTransactionByIDRequest) (res model.GetTransactionResponse, err error)
}
