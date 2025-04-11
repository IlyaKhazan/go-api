package storage

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func GetDBConnect(connStr string) (*pgx.Conn, error) {
	slog.Info("Connecting to Postgres", "url", connStr)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		return nil, errors.Wrap(err, "unable to connect to database")
	}

	slog.Info("Connected to PostgreSQL successfully")
	return conn, nil
}
