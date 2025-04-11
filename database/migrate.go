package database

import (
	"database/sql"
	"embed"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(url string) error {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return errors.Wrap(err, "cannot connect to db")
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			slog.Warn("failed to close db connection after migration", slog.Any("error", cerr))
		}
	}()

	if err = db.Ping(); err != nil {
		return errors.Wrap(err, "cannot ping db")
	}

	goose.SetBaseFS(migrations)
	if err = goose.SetDialect("postgres"); err != nil {
		return errors.Wrap(err, "cannot set database dialect")
	}

	slog.Info("applying migrations...")

	if err = goose.Up(db, "migrations"); err != nil {
		return errors.Wrap(err, "cannot apply migrations")
	}

	slog.Info("migrations applied successfully")
	return nil
}
