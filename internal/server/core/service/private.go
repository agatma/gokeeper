package service

import (
	"context"
	"errors"
	"fmt"
	"gokeeper/internal/server/adapters/storage"
	domain2 "gokeeper/pkg/domain"

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

func (ps *PrivateService) Save(ctx context.Context, pd *domain2.Data, userID uuid.UUID) error {
	tx, err := ps.privateStorage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	existingPrivateData, err := ps.privateStorage.GetByID(ctx, pd.ID, userID, tx)
	if err != nil && !errors.Is(err, domain2.ErrPrivateDataNotFound) {
		return fmt.Errorf("failed to get existing private data: %w", err)
	}

	if existingPrivateData != nil && existingPrivateData.SavedAt.After(pd.SavedAt) {
		return domain2.ErrPrivateDataConflict
	}

	if err = ps.privateStorage.InsertOrUpdate(ctx, pd, userID, tx); err != nil {
		return fmt.Errorf("failed to save private data: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (ps *PrivateService) GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain2.Data, error) {
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

func (ps *PrivateService) Delete(ctx context.Context, pd *domain2.DeleteRequest, userID uuid.UUID) error {
	tx, err := ps.privateStorage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	existingPrivateData, err := ps.privateStorage.GetByID(ctx, pd.ID, userID, tx)
	if err != nil {
		switch {
		case errors.Is(err, domain2.ErrPrivateDataNotFound):
			return nil
		}
		return fmt.Errorf("failed to get existing private data: %w", err)
	}

	if existingPrivateData != nil && existingPrivateData.SavedAt.After(pd.DeletedAt) {
		return domain2.ErrPrivateDataConflict
	}

	if err = ps.privateStorage.Delete(ctx, pd.ID, userID, tx); err != nil {
		return fmt.Errorf("failed to delete private data: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (ps *PrivateService) GetAll(ctx context.Context, req *domain2.GetAllRequest, userID uuid.UUID) ([]domain2.Data, error) {
	data, err := ps.privateStorage.GetAll(ctx, req, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get private data: %w", err)
	}
	return data, nil
}
