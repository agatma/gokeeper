package postgresql

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"gokeeper/pkg/logger"

	"github.com/pressly/goose/v3"
)

//go:embed migrations
var migrations embed.FS

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("postgres migrate set dialect postgres: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		if !errors.Is(err, goose.ErrNoNextVersion) {
			return fmt.Errorf("postgres migrate up: %w", err)
		}
	}
	logger.Log.Info("successful migrations")
	return nil
}
