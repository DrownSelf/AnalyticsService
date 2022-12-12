package controllers

import (
	"context"
	"log"

	"github.com/DrownSelf/AnalyticsService/internal/entities"
	"github.com/DrownSelf/AnalyticsService/pkg/provider"
)

type AnalyticsRepository interface {
	AddDriverLog(context.Context, entities.DriverLog) error
	AddUserLog(context.Context, entities.UserLog) error
	AddOrderLog(context.Context, entities.OrderLog) error
	GetDriverRegisterStats(context.Context) ([]entities.DriverStatistic, error)
	GetUserRegisterStats(context.Context) ([]entities.UserStatistic, error)
	GetOrderStats(ctx context.Context) ([]entities.OrderStatistic, error)
	CheckAdminAccount(ctx context.Context, username string, password string) error
}

type AnalyticsService struct {
	driverConsumer      *provider.Consumer[entities.DriverLog]
	orderConsumer       *provider.Consumer[entities.OrderLog]
	userConsumer        *provider.Consumer[entities.UserLog]
	analyticsRepository AnalyticsRepository
}

func NewAnalyticsService(driverConsumer *provider.Consumer[entities.DriverLog], orderConsumer *provider.Consumer[entities.OrderLog], userConsumer *provider.Consumer[entities.UserLog], repository AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{driverConsumer: driverConsumer, orderConsumer: orderConsumer, userConsumer: userConsumer, analyticsRepository: repository}
}

func (a *AnalyticsService) AddDriverLog(ctx context.Context) {
	a.driverConsumer.Read(entities.DriverLog{}, func(driverLog entities.DriverLog, err error) {
		if err != nil {
			log.Println("can't insert data because of error in kafka")
			return
		}
		if err = a.analyticsRepository.AddDriverLog(ctx, driverLog); err != nil {
			log.Printf("repository error: %v", err)
			return
		}
	})
}

func (a *AnalyticsService) AddUserLog(ctx context.Context) {
	a.userConsumer.Read(entities.UserLog{}, func(userLog entities.UserLog, err error) {
		if err != nil {
			log.Println("can't insert data because of error in kafka")
			return
		}
		if err = a.analyticsRepository.AddUserLog(ctx, userLog); err != nil {
			log.Printf("repository error: %v", err)
			return
		}
	})
}

func (a *AnalyticsService) AddOrderLog(ctx context.Context) {
	a.orderConsumer.Read(entities.OrderLog{}, func(orderLog entities.OrderLog, err error) {
		if err != nil {
			log.Println("can't insert data because of error in kafka")
			return
		}
		if err = a.analyticsRepository.AddOrderLog(ctx, orderLog); err != nil {
			log.Printf("repository error: %v", err)
			return
		}
	})
}

func (a *AnalyticsService) GetUserRegisterStats(ctx context.Context) ([]entities.UserStatistic, error) {
	stats, err := a.analyticsRepository.GetUserRegisterStats(ctx)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (a *AnalyticsService) GetDriverRegisterStats(ctx context.Context) ([]entities.DriverStatistic, error) {
	stats, err := a.analyticsRepository.GetDriverRegisterStats(ctx)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (a *AnalyticsService) GetOrderStats(ctx context.Context) ([]entities.OrderStatistic, error) {
	stats, err := a.analyticsRepository.GetOrderStats(ctx)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (a *AnalyticsService) LogIn(ctx context.Context, username string, password string) error {
	if err := a.analyticsRepository.CheckAdminAccount(ctx, username, password); err != nil {
		return err
	}
	return nil
}
