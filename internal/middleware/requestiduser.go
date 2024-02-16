package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func RequestIdUser() fiber.Handler {
	return func(c *fiber.Ctx) error {

		val, ok := c.Locals(requestid.ConfigDefault.ContextKey).(string)
		if ok {
			userCtx := c.UserContext()
			userCtx = context.WithValue(userCtx, requestid.ConfigDefault.ContextKey, val)
			c.SetUserContext(userCtx)
		}

		return c.Next()
	}
}
