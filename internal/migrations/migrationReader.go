package migrations

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type ClickHouseMigrator struct {
	migrator *migrate.Migrate
}

func ConnectMigrator(db *sql.DB) (ClickHouseMigrator, error) {
	driver, err := clickhouse.WithInstance(db, &clickhouse.Config{})
	if err != nil {
		return ClickHouseMigrator{}, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://./internal/migrations", "analytics", driver)
	if err != nil {
		return ClickHouseMigrator{}, err
	}

	return ClickHouseMigrator{m}, nil
}

func (c *ClickHouseMigrator) Up() error {
	if err := c.migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}
	return nil
}

func (c *ClickHouseMigrator) Down() error {
	if err := c.migrator.Down(); err != nil {
		return err
	}
	return nil
}
