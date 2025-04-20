package auth

import (
	"context"
	"gokeeper/pkg/domain"
	"log"
)

type Service struct {
	jwtFileWorker JwtFileWorker
	authClient    Client
}
type JwtFileWorker interface {
	Set(jwt string) error
	Get() (string, error)
}

type Client interface {
	Login(ctx context.Context, user domain.InUserRequest) (string, error)
	Register(ctx context.Context, user domain.InUserRequest) (string, error)
}

func NewAuthService(jwtFileWorker JwtFileWorker, authClient Client) *Service {
	return &Service{
		jwtFileWorker: jwtFileWorker,
		authClient:    authClient,
	}
}

func (as *Service) Register(ctx context.Context, user domain.InUserRequest, saveJWT bool) error {
	jwt, err := as.authClient.Register(ctx, user)
	if err != nil {
		return err
	}

	if saveJWT && jwt != "" {
		err = as.jwtFileWorker.Set(jwt)
		if err != nil {
			log.Printf("Warn: Не удалось сохранить токен на диск: %v", err)
		}
	}
	return nil
}

func (as *Service) Login(ctx context.Context, user domain.InUserRequest, saveJWT bool) (string, error) {
	jwt, err := as.authClient.Login(ctx, user)
	if err != nil {
		return "", err
	}

	if saveJWT && jwt != "" {
		err = as.jwtFileWorker.Set(jwt)
		if err != nil {
			log.Printf("Warn: %v", err)
		}
	}
	return jwt, nil
}

func (as *Service) GetJwt(ctx context.Context) (string, error) {
	return as.jwtFileWorker.Get()
}
