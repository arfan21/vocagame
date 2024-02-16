package exception

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/arfan21/vocagame/pkg/constant"
	"github.com/arfan21/vocagame/pkg/logger"
	"github.com/arfan21/vocagame/pkg/pkgutil"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func FiberErrorHandler(ctx *fiber.Ctx, err error) error {
	defer func() {
		logger.Log(ctx.UserContext()).Error().Msg(err.Error())
	}()

	defaultRes := pkgutil.HTTPResponse{
		Code:    fiber.StatusInternalServerError,
		Message: "Internal Server Error",
	}

	var errValidation *constant.ErrValidation
	if errors.As(err, &errValidation) {
		data := errValidation.Error()
		var messages []map[string]interface{}

		errJson := json.Unmarshal([]byte(data), &messages)
		PanicIfNeeded(errJson)

		defaultRes.Code = fiber.StatusBadRequest
		defaultRes.Message = "Bad Request"
		var errors []interface{}
		for _, message := range messages {
			errors = append(errors, message)
		}
		defaultRes.Errors = errors
	}

	var notFoundError *constant.ErrNotFound
	if errors.As(err, &notFoundError) {
		defaultRes.Code = fiber.StatusNotFound
		defaultRes.Message = "Not Found"
	}

	var unauthorizedError *constant.ErrUnauthorized
	if errors.As(err, &unauthorizedError) {
		defaultRes.Code = fiber.StatusUnauthorized
		defaultRes.Message = "Unauthorized"
	}

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		defaultRes.Code = fiberError.Code
		defaultRes.Message = fiberError.Message
	}

	if errors.Is(err, pgx.ErrNoRows) {
		defaultRes.Code = fiber.StatusNotFound
		defaultRes.Message = "data not found"
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		defaultRes.Code = fiber.StatusUnprocessableEntity
		defaultRes.Message = http.StatusText(fiber.StatusUnprocessableEntity)

		defaultRes.Errors = []interface{}{
			map[string]interface{}{
				"field":   unmarshalTypeError.Field,
				"message": fmt.Sprintf("%s harus %s", unmarshalTypeError.Field, unmarshalTypeError.Type),
			},
		}
	}

	if errors.Is(err, constant.ErrEmailAlreadyRegistered) {
		defaultRes.Code = fiber.StatusConflict
		defaultRes.Message = constant.ErrEmailAlreadyRegistered.Error()
	}

	if errors.Is(err, constant.ErrEmailOrPasswordInvalid) {
		defaultRes.Code = fiber.StatusUnauthorized
		defaultRes.Message = constant.ErrEmailOrPasswordInvalid.Error()
	}

	if errors.Is(err, constant.ErrUnauthorizedAccess) {
		defaultRes.Code = fiber.StatusUnauthorized
		defaultRes.Message = constant.ErrUnauthorizedAccess.Error()
	}

	if errors.Is(err, constant.ErrCategoryNotFound) {
		defaultRes.Code = fiber.StatusNotFound
		defaultRes.Message = constant.ErrCategoryNotFound.Error()
	}

	// check decimal error decode
	if strings.Contains(err.Error(), "error decoding string") &&
		strings.Contains(err.Error(), "to decimal") {
		defaultRes.Code = fiber.StatusBadRequest
		defaultRes.Message = constant.ErrStringNotDecimal.Error()
	}

	// handle error parse uuid
	if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("invalid UUID")) {
		defaultRes.Code = fiber.StatusBadRequest
		defaultRes.Message = constant.ErrInvalidUUID.Error()
	}

	if errors.Is(err, constant.ErrProductNotFound) {
		defaultRes.Code = fiber.StatusNotFound
		defaultRes.Message = constant.ErrProductNotFound.Error()
	}

	if errors.Is(err, constant.ErrProductAlreadyAddedToCart) {
		defaultRes.Code = fiber.StatusConflict
		defaultRes.Message = constant.ErrProductAlreadyAddedToCart.Error()
	}

	if errors.Is(err, constant.ErrCannotAddOwnProductToCart) {
		defaultRes.Code = fiber.StatusConflict
		defaultRes.Message = constant.ErrCannotAddOwnProductToCart.Error()
	}

	if errors.Is(err, constant.ErrPaymentMethodNotFound) {
		defaultRes.Code = fiber.StatusNotFound
		defaultRes.Message = constant.ErrPaymentMethodNotFound.Error()
	}

	if errors.Is(err, constant.ErrNoProductInCart) {
		defaultRes.Code = fiber.StatusNotFound
		defaultRes.Message = constant.ErrNoProductInCart.Error()
	}

	if errors.Is(err, constant.ErrProductStokNotEnough) {
		defaultRes.Code = fiber.StatusBadRequest
		defaultRes.Message = constant.ErrProductStokNotEnough.Error()
	}

	if errors.Is(err, constant.ErrProductNotFoundOrStok) {
		defaultRes.Code = fiber.StatusBadRequest
		defaultRes.Message = constant.ErrProductNotFoundOrStok.Error()
	}

	if errors.Is(err, constant.ErrTransactionAlreadyPaid) {
		defaultRes.Code = fiber.StatusBadRequest
		defaultRes.Message = constant.ErrTransactionAlreadyPaid.Error()
	}

	if errors.Is(err, constant.ErrTransactionAlreadyPaidOrFailed) {
		defaultRes.Code = fiber.StatusBadRequest
		defaultRes.Message = constant.ErrTransactionAlreadyPaidOrFailed.Error()
	}

	if errors.Is(err, constant.ErrPaymentNotEqualTotalAmount) {
		defaultRes.Code = fiber.StatusBadRequest
		defaultRes.Message = constant.ErrPaymentNotEqualTotalAmount.Error()

	}

	if defaultRes.Code >= 500 {
		defaultRes.Message = http.StatusText(defaultRes.Code)
	}

	return ctx.Status(defaultRes.Code).JSON(defaultRes)
}
