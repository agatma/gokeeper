package storage

import (
	"context"
	"gokeeper/internal/server/adapters/storage/database"
	"gokeeper/internal/server/core/domain"
)

type AuthStorage interface {
	GetUser(ctx context.Context, login string) (domain.User, error)
	InsertUser(ctx context.Context, newUser domain.User, trx *database.Trx) error
	BeginTx(ctx context.Context) (*database.Trx, error)
}

type Storage interface {
	AuthStorage
}

func NewStorage(dsn string) (Storage, error) {
	return database.NewStorage(dsn)
}
