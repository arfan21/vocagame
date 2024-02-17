package transaction

import (
	"context"

	"github.com/arfan21/vocagame/internal/model"
)

type Service interface {
	CreateDepositTransaction(ctx context.Context, req model.CreateDepositTransactionRequest) (res model.CreateTransactionResponse, err error)
}
