package controllers

import (
	"context"
	"time"

	"github.com/DrownSelf/AnalyticsService/internal/entities"
)

type AnalystRepository interface {
	AddDriverLog(context.Context, time.Time) error
	AddUserLog(context.Context, time.Time) error
	AddOrderLog(context.Context, time.Time, string) error
	GetDriverRegisterStats(context.Context) ([]entities.DriverStatistic, error)
	GetUserRegisterStats(context.Context) ([]entities.UserStatistic, error)
	GetOrderStats(ctx context.Context) ([]entities.OrderStatistic, error)
}
