package storage

import (
	"context"
	"gokeeper/internal/server/adapters/storage/database"
	"gokeeper/internal/server/adapters/storage/database/postgresql"
	"gokeeper/internal/server/core/domain"

	"github.com/google/uuid"
)

type AuthStorage interface {
	GetUser(ctx context.Context, login string) (domain.User, error)
	InsertUser(ctx context.Context, newUser domain.User, tx *database.Trx) error
	BeginTx(ctx context.Context) (*database.Trx, error)
}

type PrivateStorage interface {
	GetByID(ctx context.Context, id string, userID uuid.UUID, tx *database.Trx) (*domain.Data, error)
	InsertOrUpdate(ctx context.Context, pd *domain.Data, userID uuid.UUID, tx *database.Trx) error
	Delete(ctx context.Context, id string, userID uuid.UUID, tx *database.Trx) error
	GetAll(ctx context.Context, req *domain.GetAllRequest, userID uuid.UUID) ([]domain.Data, error)
	BeginTx(ctx context.Context) (*database.Trx, error)
}

type Storage interface {
	AuthStorage
	PrivateStorage
}

func NewStorage(dsn string) (Storage, error) {
	return postgresql.NewStorage(dsn)
}
