package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"gokeeper/internal/server/adapters/storage"
	"gokeeper/pkg/auth"
	domain2 "gokeeper/pkg/domain"

	"github.com/google/uuid"
)

type AuthService struct {
	authStorage   storage.AuthStorage
	authenticator auth.Authenticator
}

func NewAuthService(authStorage storage.AuthStorage, authenticator auth.Authenticator) *AuthService {
	return &AuthService{
		authStorage:   authStorage,
		authenticator: authenticator,
	}
}

func (as *AuthService) Register(ctx context.Context, inUser domain2.InUserRequest) (auth.Token, error) {
	tx, err := as.authStorage.BeginTx(ctx)

	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %v", err)
	}

	if registeredUser, err := as.authStorage.GetUser(ctx, inUser.Login); err == nil {
		if registeredUser.Login != "" {
			return "", domain2.ErrUserConflict
		}
	}
	newUser := domain2.User{
		Login:        inUser.Login,
		PasswordHash: generatePasswordHash(inUser.Password),
		ID:           uuid.New(),
	}
	err = as.authStorage.InsertUser(ctx, newUser, tx)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return "", fmt.Errorf("failed to rollback transaction: %w", err)
		}
		return "", fmt.Errorf("failed to insert user: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	token, err := as.authenticator.MakeJWT(newUser.ID, newUser.Login)
	if err != nil {
		return "", errors.New("failed to generate token, %")
	}
	return token, nil
}

func (as *AuthService) Login(ctx context.Context, inUser domain2.InUserRequest) (auth.Token, error) {
	tx, err := as.authStorage.BeginTx(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %v", err)
	}

	userInDB, err := as.authStorage.GetUser(ctx, inUser.Login)
	if err != nil {
		if errors.Is(err, domain2.ErrUserNotFound) {
			return "", domain2.ErrUserAuthentication
		}
		return "", err
	}

	if !checkPassword(userInDB.PasswordHash, inUser.Password) {
		return "", domain2.ErrUserAuthentication
	}
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	token, err := as.authenticator.MakeJWT(userInDB.ID, userInDB.Login)
	if err != nil {
		return "", errors.New("failed to generate token, %")
	}
	return token, nil
}

func generatePasswordHash(password string) []byte {
	h := sha256.New()
	h.Write([]byte(password))
	hash := h.Sum(nil)
	return hash
}

func checkPassword(passHash []byte, password string) bool {
	return bytes.Equal(passHash, generatePasswordHash(password))
}
