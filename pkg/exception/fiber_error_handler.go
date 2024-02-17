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

	var badRequestError *constant.ErrNotFound
	if errors.As(err, &badRequestError) {
		defaultRes.Code = fiber.StatusBadRequest
		if badRequestError.Message != "" {
			defaultRes.Message = badRequestError.Message
		} else {
			defaultRes.Message = "Bad Request"
		}
	}

	var notFoundError *constant.ErrNotFound
	if errors.As(err, &notFoundError) {
		defaultRes.Code = fiber.StatusNotFound
		if notFoundError.Message != "" {
			defaultRes.Message = notFoundError.Message
		} else {
			defaultRes.Message = "Not Found"
		}
	}

	var unauthorizedError *constant.ErrUnauthorized
	if errors.As(err, &unauthorizedError) {
		defaultRes.Code = fiber.StatusUnauthorized
		if unauthorizedError.Message != "" {
			defaultRes.Message = unauthorizedError.Message
		} else {
			defaultRes.Message = "Unauthorized"
		}
	}

	var forbiddenError *constant.ErrForbidden
	if errors.As(err, &forbiddenError) {
		defaultRes.Code = fiber.StatusForbidden
		if forbiddenError.Message != "" {
			defaultRes.Message = forbiddenError.Message
		} else {
			defaultRes.Message = "Forbidden"
		}
	}

	var conflictError *constant.ErrConflict
	if errors.As(err, &conflictError) {
		defaultRes.Code = fiber.StatusConflict
		if conflictError.Message != "" {
			defaultRes.Message = conflictError.Message
		} else {
			defaultRes.Message = "Conflict"
		}
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

	if defaultRes.Code >= 500 {
		defaultRes.Message = http.StatusText(defaultRes.Code)
	}

	return ctx.Status(defaultRes.Code).JSON(defaultRes)
}
