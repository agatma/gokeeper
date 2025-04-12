package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gokeeper/internal/server/core/domain"
)

type Storage struct {
	db  *sql.DB
	dsn string
}

func NewStorage(dsn string) (*Storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open pg connection: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database %w", err)
	}
	if err = Migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database %w", err)
	}
	return &Storage{
		dsn: dsn,
		db:  db,
	}, nil
}

func (s Storage) GetUser(ctx context.Context, login string) (domain.User, error) {
	getUserFromDB := `
		SELECT id, login, password_hash FROM users WHERE login = $1;
	`
	row := s.db.QueryRowContext(ctx, getUserFromDB, login)

	var userInDB domain.User
	err := row.Scan(&userInDB.ID, &userInDB.Login, &userInDB.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return userInDB, nil
}

func (s Storage) InsertUser(ctx context.Context, newUser domain.User, tx *Trx) error {
	insertUserQuery := `
		INSERT INTO users (id, login, password_hash) VALUES ($1, $2, $3);
	`
	_, err := tx.ExecContext(ctx, insertUserQuery, newUser.ID, newUser.Login, newUser.PasswordHash)
	return err
}

func (s Storage) BeginTx(ctx context.Context) (*Trx, error) {
	return BeginTx(ctx, s.db)
}
