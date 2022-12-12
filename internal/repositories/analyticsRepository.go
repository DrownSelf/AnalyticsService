package repositories

import (
	"context"
	"crypto/tls"
	"database/sql"
	"net"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/DrownSelf/AnalyticsService/internal/configs"
	"github.com/DrownSelf/AnalyticsService/internal/entities"
	"github.com/DrownSelf/AnalyticsService/internal/migrations"
)

type AnalyticsRepository struct {
	clickhouse *sql.DB
}

func NewAnalyticsRepo(config *configs.Config) (*AnalyticsRepository, error) {
	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{config.ClickHouseConnection},
		Auth: clickhouse.Auth{
			Database: config.ClickHouseDbName,
			Username: config.ClickHouseUserName,
			Password: config.ClickHousePassword,
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
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

func (a *AnalyticsRepository) AddOrderLog(ctx context.Context, log entities.OrderLog) error {
	query := `insert into orders(place, orderDate) values ($1, $2);`
	if _, err := a.clickhouse.ExecContext(ctx, query, log.OrderTimeLog, log.EndPlace); err != nil {
		return err
	}
	return nil
}

func (a *AnalyticsRepository) AddUserLog(ctx context.Context, log entities.UserLog) error {
	query := `insert into users(registrationDate) values ($1)`
	if _, err := a.clickhouse.ExecContext(ctx, query, log.RegistrationLog); err != nil {
		return err
	}
	return nil
}

func (a *AnalyticsRepository) AddDriverLog(ctx context.Context, log entities.DriverLog) error {
	query := `insert into drivers(registrationDate) values ($1)`
	if _, err := a.clickhouse.ExecContext(ctx, query, log.RegistrationLog); err != nil {
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

func (a *AnalyticsRepository) CheckAdminAccount(ctx context.Context, username string, password string) error {
	query := `select * from analyst where password = $1 and username = $2`
	_, err := a.clickhouse.QueryContext(ctx, query, password, username)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fiber.NewError(fiber.StatusBadRequest, "Wrong username or password")
		}
		return err
	}
	return nil
}
