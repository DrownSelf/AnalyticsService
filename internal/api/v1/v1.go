package v1

import (
	"github.com/gofiber/fiber/v2"

	"github.com/DrownSelf/AnalyticsService/internal/handlers"
	"github.com/DrownSelf/AnalyticsService/internal/middlewares"
)

type ApiV1 struct {
	handler *handlers.Handler
}

func NewApiV1(handler *handlers.Handler) *ApiV1 {
	return &ApiV1{handler: handler}
}

func (a *ApiV1) InitApiV1Groups(router fiber.Router) {
	v1 := router.Group("/v1")
	v1.Post("/login")
	v1.Get("/logout")

	analystGroup := v1.Group("/analyst", middlewares.Authorize)
	analystGroup.Get("/ordersStats")
	analystGroup.Get("/userStats")
	analystGroup.Get("/driverStats")
}
