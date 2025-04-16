package service

import (
	"context"
	"errors"
	"fmt"
	"gokeeper/internal/server/adapters/storage"
	"gokeeper/internal/server/core/domain"

	"github.com/google/uuid"
)

type PrivateService struct {
	privateStorage storage.PrivateStorage
}

func NewPrivateService(privateStorage storage.PrivateStorage) *PrivateService {
	return &PrivateService{
		privateStorage: privateStorage,
	}
}

func (ps *PrivateService) Save(ctx context.Context, pd *domain.Data, userID uuid.UUID) error {
	tx, err := ps.privateStorage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	existingPrivateData, err := ps.privateStorage.GetByID(ctx, pd.ID, userID, tx)
	if err != nil && !errors.Is(err, domain.ErrPrivateDataNotFound) {
		return fmt.Errorf("failed to get existing private data: %w", err)
	}

	if existingPrivateData != nil && existingPrivateData.SavedAt.After(pd.SavedAt) {
		return domain.ErrPrivateDataConflict
	}

	if err = ps.privateStorage.InsertOrUpdate(ctx, pd, userID, tx); err != nil {
		return fmt.Errorf("failed to save private data: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (ps *PrivateService) GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain.Data, error) {
	tx, err := ps.privateStorage.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	existingPrivateData, err := ps.privateStorage.GetByID(ctx, id, userID, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing private data: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return existingPrivateData, nil
}

func (ps *PrivateService) Delete(ctx context.Context, pd *domain.DeleteRequest, userID uuid.UUID) error {
	tx, err := ps.privateStorage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	existingPrivateData, err := ps.privateStorage.GetByID(ctx, pd.ID, userID, tx)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPrivateDataNotFound):
			return nil
		}
		return fmt.Errorf("failed to get existing private data: %w", err)
	}

	if existingPrivateData != nil && existingPrivateData.SavedAt.After(pd.DeletedAt) {
		return domain.ErrPrivateDataConflict
	}

	if err = ps.privateStorage.Delete(ctx, pd.ID, userID, tx); err != nil {
		return fmt.Errorf("failed to delete private data: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (ps *PrivateService) GetAll(ctx context.Context, req *domain.GetAllRequest, userID uuid.UUID) ([]domain.Data, error) {
	data, err := ps.privateStorage.GetAll(ctx, req, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get private data: %w", err)
	}
	return data, nil
}
