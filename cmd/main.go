package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"

	"github.com/DrownSelf/AnalyticsService/internal/api"
	"github.com/DrownSelf/AnalyticsService/internal/api/v1"
	"github.com/DrownSelf/AnalyticsService/internal/appErrors"
	"github.com/DrownSelf/AnalyticsService/internal/configs"
	"github.com/DrownSelf/AnalyticsService/internal/controllers"
	"github.com/DrownSelf/AnalyticsService/internal/entities"
	"github.com/DrownSelf/AnalyticsService/internal/handlers"
	"github.com/DrownSelf/AnalyticsService/internal/repositories"
	"github.com/DrownSelf/AnalyticsService/pkg/provider"
)

func main() {
	config, err := configs.LoadConnectionConfig()
	if err != nil {
		log.Fatalf("error during extract config: %v", err)
	}

	repository, err := repositories.NewAnalyticsRepo(config)
	if err != nil {
		log.Fatalf("error during connect to db: %v", err)
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	driverConsumer := provider.Consumer[entities.DriverLog]{}
	if err = driverConsumer.CreateConnection([]string{config.ConsumerConnection}, dialer, ""); err != nil {
		log.Fatalf("error during connect to message broker: %v", err)
	}

	orderConsumer := provider.Consumer[entities.OrderLog]{}
	if err = orderConsumer.CreateConnection([]string{config.ConsumerConnection}, dialer, ""); err != nil {
		log.Fatalf("error during connect to message broker: %v", err)
	}

	userConsumer := provider.Consumer[entities.UserLog]{}
	if err = userConsumer.CreateConnection([]string{config.ConsumerConnection}, dialer, ""); err != nil {
		log.Fatalf("error during connect to message broker: %v", err)
	}

	service := controllers.NewAnalyticsService(&driverConsumer, &orderConsumer, &userConsumer, repository)
	handler := handlers.NewHandler(service)
	app := fiber.New(fiber.Config{
		ErrorHandler: appErrors.HandleError,
	})
	v1 := v1.NewApiV1(handler)
	api := api.NewApiGroups(v1)
	api.InitRouterGroups(app)
	go func() {
		log.Fatal(app.Listen(config.ApiPort))
	}()
	go func() {
		service.AddUserLog(context.Background())
	}()
	go func() {
		service.AddDriverLog(context.Background())
	}()
	go func() {
		service.AddOrderLog(context.Background())
	}()
}
