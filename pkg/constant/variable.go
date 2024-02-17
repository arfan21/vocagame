package constant

type ContextKey string

var (
	JWTClaimsContextKey ContextKey
)

const (
	TransactionTypeDepositID  = 1
	TransactionTypeWithdrawID = 2
	TransactionTypePurchaseID = 3
	TransactionTypeRefundID   = 4
)
