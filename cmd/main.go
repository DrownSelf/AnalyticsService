package main

import (
	"context"
	"log"
	"time"

	"github.com/DrownSelf/AnalyticsService/internal/repositories"
)

func main() {
	repo, err := repositories.NewAnalyticsRepo()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	stats, err := repo.GetOrderStats(context.Background())
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Print(stats)
	_ = repo.AddDriverLog(context.TODO(), time.Now().AddDate(0, 2, 14))
	statsOfDriver, err := repo.GetDriverRegisterStats(context.Background())
	log.Print(statsOfDriver)
}
