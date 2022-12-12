package appErrors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func HandleError(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return ctx.Status(code).JSON(err.Error())
}
