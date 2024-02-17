package transactionctrl

import (
	"github.com/arfan21/vocagame/internal/model"
	"github.com/arfan21/vocagame/internal/transaction"
	"github.com/arfan21/vocagame/pkg/constant"
	"github.com/arfan21/vocagame/pkg/exception"
	"github.com/arfan21/vocagame/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ControllerHTTP struct {
	svc transaction.Service
}

func New(svc transaction.Service) *ControllerHTTP {
	return &ControllerHTTP{svc: svc}
}

// @Summary Create Deposit Transaction
// @Description Create Deposit Transaction
// @Tags Transaction
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.CreateDepositTransactionRequest true "Create Deposit Transaction"
// @Success 201 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{errors=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /api/v1/transactions/deposit [post]
func (ctrl ControllerHTTP) CreateDepositTransaction(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(pkgutil.HTTPResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "invalid or expired token",
		})
	}

	var req model.CreateDepositTransactionRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	uuidUserID, err := uuid.Parse(claims.Subject)
	exception.PanicIfNeeded(err)
	req.UserID = uuidUserID

	id, err := ctrl.svc.CreateDepositTransaction(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusCreated).JSON(pkgutil.HTTPResponse{
		Code: fiber.StatusCreated,
		Data: id,
	})
}