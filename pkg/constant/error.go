package constant

import "errors"

const (
	ErrSQLUniqueViolation = "23505"
)

var (
	ErrEmailAlreadyRegistered         = errors.New("email already registered")
	ErrEmailOrPasswordInvalid         = errors.New("email or password invalid")
	ErrUnauthorizedAccess             = errors.New("unauthorized access")
	ErrStringNotDecimal               = errors.New("string value is not decimal")
	ErrInvalidUUID                    = errors.New("invalid uuid length or format")
	ErrProductNotFound                = errors.New("product not found")
	ErrProductStokNotEnough           = errors.New("product stok not enough")
	ErrProductNotFoundOrStok          = errors.New("product not found or stok not enough")
	ErrTransactionAlreadyPaid         = errors.New("transaction already paid")
	ErrTransactionAlreadyPaidOrFailed = errors.New("transaction already paid or failed")
	ErrTxDetailInsertedNotEqual       = errors.New("transaction detail inserted not equal with transaction detail request")
	ErrCannotUpdateNotOwner           = &ErrForbidden{Message: "cannot update product, not owner"}
)

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
