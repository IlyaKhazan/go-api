package migrations

import (
	"database/sql"
	"embed"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migrations embed.FS

func Migrate(url string) error {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return errors.Wrap(err, "cannot connect to db")
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	if err = db.Ping(); err != nil {
		return errors.Wrap(err, "cannot ping db")
	}

	goose.SetBaseFS(migrations)
	if err = goose.SetDialect("postgres"); err != nil {
		return errors.Wrap(err, "cannot set migrations dialect")
	}

	version, err := goose.GetDBVersion(db)
	if err != nil {
		return errors.Wrap(err, "cannot get migration version")
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		if err = goose.DownTo(db, "migrations", version); err != nil {
			slog.Error(
				"cannot rollback migrations",
				slog.Any("error", err),
				slog.Any("try rollback to version", version),
			)
		}

		return errors.Wrap(err, "cannot up migrations")
	}

	return nil
}
