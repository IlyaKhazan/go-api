package storage

import (
	"context"
	"fmt"
	"log/slog"

	"go-api/config"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func GetDBConnect(cfg *config.Config) (*pgx.Conn, error) {
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
