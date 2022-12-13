package api

import (
	"github.com/gofiber/fiber/v2"

	"github.com/DrownSelf/AnalyticsService/internal/api/v1"
)

type ApiGroup struct {
	api *v1.ApiV1
}

func NewApiGroups(api *v1.ApiV1) *ApiGroup {
	return &ApiGroup{api: api}
}

func (a *ApiGroup) InitRouterGroups(app *fiber.App) {
	api := app.Group("/api")
	a.api.InitApiV1Groups(api)
}
