package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/DrownSelf/AnalyticsService/internal/entities"
)

type AnalyticsService interface {
	AddDriverLog(context.Context)
	AddUserLog(context.Context)
	AddOrderLog(context.Context)
	GetDriverRegisterStats(context.Context) ([]entities.DriverStatistic, error)
	GetUserRegisterStats(context.Context) ([]entities.UserStatistic, error)
	GetOrderStats(context.Context) ([]entities.OrderStatistic, error)
	LogIn(ctx context.Context, username string, password string) error
}

type Handler struct {
	service AnalyticsService
}

func (h *Handler) LogIn(ctx *fiber.Ctx) error {
	var request LogInRequest
	if err := ctx.BodyParser(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid data")
	}

	if err := h.service.LogIn(ctx.Context(), request.Username, request.Password); err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:   request.Username,
		Value:  request.Password,
		MaxAge: 3600,
	})
	_ = ctx.Status(fiber.StatusOK).JSON("Logged in successfully")
	return nil
}

func (h *Handler) LogOut(ctx *fiber.Ctx) error {
	name := ctx.Cookies("Name")
	if name == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "you even not authorized")
	}
	ctx.ClearCookie(name)
	_ = ctx.Status(fiber.StatusOK).JSON("Logged out successfully")
	return nil
}

func (h *Handler) GetDriverStats(ctx *fiber.Ctx) error {
	driverStats, err := h.service.GetDriverRegisterStats(ctx.Context())
	if err != nil {
		return err
	}
	_ = ctx.Status(fiber.StatusOK).JSON(driverStats)
	return nil
}

func (h *Handler) GetUserStats(ctx *fiber.Ctx) error {
	userStats, err := h.service.GetUserRegisterStats(ctx.Context())
	if err != nil {
		return err
	}
	_ = ctx.Status(fiber.StatusOK).JSON(userStats)
	return nil
}

func (h *Handler) GetOrderStats(ctx *fiber.Ctx) error {
	orderStats, err := h.service.GetOrderStats(ctx.Context())
	if err != nil {
		return err
	}
	_ = ctx.Status(fiber.StatusOK).JSON(orderStats)
	return nil
}
