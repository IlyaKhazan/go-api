package config

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Address    string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Enable environment variables
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("No config file found, using system environment variables")
	}

	config := &Config{
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		Address:    viper.GetString("ADDRESS"),
	}

	slog.Info("Configuration loaded successfully")
	return config, nil
}

func GetDBConnect(cfg *Config) (*pgx.Conn, error) {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	slog.Info("Connecting to Postgres", "url", dbURL)

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		return nil, errors.Wrap(err, "unable to connect to database")
	}

	slog.Info("Connected to PostgreSQL successfully")
	return conn, nil
}
