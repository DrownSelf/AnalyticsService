package repositories

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"

	"github.com/DrownSelf/AnalyticsService/internal/entities"
	"github.com/DrownSelf/AnalyticsService/internal/migrations"
)

type AnalyticsRepository struct {
	clickhouse *sql.DB
}

func NewAnalyticsRepo() (*AnalyticsRepository, error) {
	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "analytics",
			Username: "drown",
			Password: "150869",
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: true,
		Debugf: func(format string, v ...interface{}) {
			fmt.Println(format, v)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout: time.Duration(10) * time.Second,
	})
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)

	migrator, err := migrations.ConnectMigrator(db)
	if err != nil {
		return nil, err
	}

	if err = migrator.Up(); err != nil {
		return nil, err
	}

	return &AnalyticsRepository{db}, nil
}

func (a *AnalyticsRepository) DestroyRepository(ctx context.Context) error {
	errorChan := make(chan error)
	go func(chan error) {
		errorChan <- a.clickhouse.Close()
	}(errorChan)

	select {
	case <-ctx.Done():
		return nil
	case err := <-errorChan:
		return err
	}
}

func (a *AnalyticsRepository) AddOrderLog(ctx context.Context, timeOfOrder time.Time, region string) error {
	query := `insert into orders(place, orderDate) values ($1, $2);`
	if _, err := a.clickhouse.ExecContext(ctx, query, timeOfOrder, region); err != nil {
		return err
	}
	return nil
}

func (a *AnalyticsRepository) AddUserLog(ctx context.Context, timeOfRegistration time.Time) error {
	query := `insert into users(registrationDate) values ($1)`
	if _, err := a.clickhouse.ExecContext(ctx, query, timeOfRegistration); err != nil {
		return err
	}
	return nil
}

func (a *AnalyticsRepository) AddDriverLog(ctx context.Context, timeOfRegistration time.Time) error {
	query := `insert into drivers(registrationDate) values ($1)`
	if _, err := a.clickhouse.ExecContext(ctx, query, timeOfRegistration); err != nil {
		return err
	}
	return nil
}

func (a *AnalyticsRepository) GetDriverRegisterStats(ctx context.Context) ([]entities.DriverStatistic, error) {
	var driverStats []entities.DriverStatistic
	query := `select MONTH(registrationDate) as month, count(month) as count from drivers group by month order by count DESC;`
	stats, err := a.clickhouse.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for stats.Next() {
		var statistic entities.DriverStatistic
		err = stats.Scan(&statistic.Month, &statistic.Count)

		if err != nil {
			return nil, err
		}
		driverStats = append(driverStats, statistic)
	}
	return driverStats, nil
}

func (a *AnalyticsRepository) GetUserRegisterStats(ctx context.Context) ([]entities.UserStatistic, error) {
	var userStats []entities.UserStatistic
	query := `select MONTH(registrationDate) as month, count(month) as count from users group by month order by count DESC;`
	stats, err := a.clickhouse.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for stats.Next() {
		var statistic entities.UserStatistic
		err = stats.Scan(statistic.Month, statistic.Count)

		if err != nil {
			return nil, err
		}
		userStats = append(userStats, statistic)
	}
	return userStats, nil
}

func (a *AnalyticsRepository) GetOrderStats(ctx context.Context) ([]entities.OrderStatistic, error) {
	var orderStats []entities.OrderStatistic
	query := `select count(place) as counter, place, HOUR(orderDate) as hour from orders group by place, hour order by counter DESC;`
	stats, err := a.clickhouse.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for stats.Next() {
		var statistic entities.OrderStatistic
		err = stats.Scan(&statistic.Count, &statistic.Place, &statistic.Hour)

		if err != nil {
			return nil, err
		}
		orderStats = append(orderStats, statistic)
	}
	return orderStats, nil
}
