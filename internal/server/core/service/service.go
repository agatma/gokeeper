package service

import (
	"gokeeper/internal/server/adapters/storage"
	"gokeeper/pkg/auth"
)

type Services struct {
	*AuthService
}

func NewServices(
	storage storage.Storage,
	authenticator auth.Authenticator,
) *Services {
	return &Services{
		NewAuthService(storage, authenticator),
	}
}
