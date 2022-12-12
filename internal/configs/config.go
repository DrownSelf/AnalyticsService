package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ClickHouseConnection string
	ConsumerConnection   string
	ClickHouseUserName   string
	ClickHousePassword   string
	ClickHouseDbName     string
}

func LoadConnectionConfig() (*Config, error) {
	var config Config
	err := godotenv.Load("./internal/config/connection.env")
	if err != nil {
		return nil, err
	}

	config.ClickHouseConnection = os.Getenv("CH_CONNECTION")
	config.ConsumerConnection = os.Getenv("KAFKA_CONNECTION")
	config.ClickHouseDbName = os.Getenv("CH_DB_NAME")
	config.ClickHouseUserName = os.Getenv("CH_USERNAME")
	config.ClickHousePassword = os.Getenv("CH_PASSWORD")
	return &config, err
}
