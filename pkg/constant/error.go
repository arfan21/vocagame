package constant

import "errors"

const (
	ErrSQLUniqueViolation = "23505"
)

var (
	ErrEmailAlreadyRegistered         = &ErrConflict{Message: "email already registered"}
	ErrEmailOrPasswordInvalid         = &ErrUnauthorized{Message: "email or password invalid"}
	ErrUnauthorizedAccess             = &ErrUnauthorized{Message: "unauthorized access"}
	ErrStringNotDecimal               = &ErrBadRequest{Message: "string not decimal"}
	ErrInvalidUUID                    = &ErrBadRequest{Message: "invalid UUID"}
	ErrProductNotFound                = &ErrNotFound{Message: "product not found"}
	ErrProductStokNotEnough           = &ErrBadRequest{Message: "product stok not enough"}
	ErrProductNotFoundOrStok          = &ErrBadRequest{Message: "product not found or stok not enough"}
	ErrTransactionAlreadyPaid         = &ErrBadRequest{Message: "transaction already paid"}
	ErrTransactionAlreadyPaidOrFailed = &ErrBadRequest{Message: "transaction already paid or failed"}
	ErrTxDetailInsertedNotEqual       = errors.New("transaction detail inserted not equal with transaction detail request")
	ErrCannotUpdateNotOwner           = &ErrForbidden{Message: "cannot update product, not owner"}
	ErrCannotDeleteNotOwner           = &ErrForbidden{Message: "cannot delete product, not owner"}
	ErrWalletAlreadyCreated           = &ErrConflict{Message: "wallet already created"}
	ErrWalletNotFound                 = &ErrNotFound{Message: "wallet not found"}
	ErrInsufficientBalance            = &ErrBadRequest{Message: "insufficient balance"}
	ErrCannotPurchaseOwnProduct       = &ErrBadRequest{Message: "cannot purchase own product"}
	ErrTransactionNotFound            = &ErrNotFound{Message: "transaction not found"}
)

type ErrBadRequest struct {
	Message string
}

func (e *ErrBadRequest) Error() string {
	return e.Message
}

type ErrNotFound struct {
	Message string
}

func (e *ErrNotFound) Error() string {
	return e.Message
}

type ErrUnauthorized struct {
	Message string
}

func (e *ErrUnauthorized) Error() string {
	return e.Message
}

type ErrValidation struct {
	Message string
}

func (e *ErrValidation) Error() string {
	return e.Message
}

type ErrForbidden struct {
	Message string
}

func (e *ErrForbidden) Error() string {
	return e.Message
}

type ErrConflict struct {
	Message string
}

func (e *ErrConflict) Error() string {
	return e.Message
}
