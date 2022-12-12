package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func Authorize(ctx *fiber.Ctx) error {
	if cookie := ctx.Cookies("admin"); cookie == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "You should authorize")
	}
	_ = ctx.Next()
	return nil
}
