package storage

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func GetDBConnect(connStr string) (*pgxpool.Pool, error) {
	slog.Info("Connecting to Postgres", "url", connStr)

	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		return nil, errors.Wrap(err, "unable to connect to database")
	}

	slog.Info("Connected to PostgreSQL successfully")
	return conn, nil
}
