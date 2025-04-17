package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gokeeper/internal/server/adapters/storage/database"
	"gokeeper/internal/server/adapters/storage/database/postgresql/queries"
	"gokeeper/pkg/domain"
	"gokeeper/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
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

func (s Storage) BeginTx(ctx context.Context) (*database.Trx, error) {
	return database.BeginTx(ctx, s.db)
}

func (s Storage) GetUser(ctx context.Context, login string) (domain.User, error) {
	row := s.db.QueryRowContext(ctx, queries.GetUser, login)

	var userInDB domain.User
	err := row.Scan(&userInDB.ID, &userInDB.Login, &userInDB.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("failed to scan user from db: %w", err)
	}
	return userInDB, nil
}

func (s Storage) InsertUser(ctx context.Context, newUser domain.User, tx *database.Trx) error {
	if _, err := tx.ExecContext(ctx, queries.InsertUser, newUser.ID, newUser.Login, newUser.PasswordHash); err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (s Storage) GetByID(ctx context.Context, id string, userID uuid.UUID, tx *database.Trx) (*domain.Data, error) {
	var privateDataInDB domain.Data
	row := tx.QueryRowContext(ctx, queries.GetDataByID, userID, id)

	privateDataInDB.ID = id
	err := row.Scan(&privateDataInDB.DataType, &privateDataInDB.Data, &privateDataInDB.MetaData, &privateDataInDB.SavedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPrivateDataNotFound
		}
		return nil, fmt.Errorf("failed to scan private data from db: %w", err)
	}
	return &privateDataInDB, nil
}

func (s Storage) InsertOrUpdate(ctx context.Context, pd *domain.Data, userID uuid.UUID, tx *database.Trx) error {
	if _, err := tx.ExecContext(ctx, queries.InsertData, pd.ID, pd.DataType, pd.Data, pd.MetaData, pd.SavedAt, userID); err != nil {
		return fmt.Errorf("failed to insert or update data: %w", err)
	}
	return nil
}

func (s Storage) Delete(ctx context.Context, id string, userID uuid.UUID, tx *database.Trx) error {
	if _, err := tx.ExecContext(ctx, queries.DeleteData, userID, id); err != nil {
		return fmt.Errorf("failed to delete data: %w", err)
	}
	return nil
}

func (s Storage) GetAll(ctx context.Context, req *domain.GetAllRequest, userID uuid.UUID) ([]domain.Data, error) {
	rows, err := s.db.QueryContext(ctx, queries.GetAllDataByUserID, userID, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			logger.Log.Error("error occurred during closing rows", zap.Error(err))
		}
	}()

	var privateData []domain.Data
	for rows.Next() {
		var privateRow domain.Data

		err = rows.Scan(&privateRow.ID, &privateRow.DataType, &privateRow.Data, &privateRow.MetaData, &privateRow.SavedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan data from db: %w", err)
		}
		privateData = append(privateData, privateRow)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate data from db: %w", err)
	}
	return privateData, nil
}
