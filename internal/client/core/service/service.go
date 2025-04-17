package service

import (
	"gokeeper/internal/client/core/service/auth"
	"gokeeper/internal/client/core/service/private"
)

type Services struct {
	AuthService    *auth.Service
	PrivateService *private.Service
}

func NewServices(
	jwtFileWorker auth.JwtFileWorker,
	authClient auth.Client,
	personalClient private.Client,
	encrypter private.Encrypter,
	privateFileWorker private.FileWorker,
	privateSender private.BulkSender,
) *Services {
	authService := auth.NewAuthService(jwtFileWorker, authClient)
	return &Services{
		AuthService:    authService,
		PrivateService: private.NewPrivateService(authService, personalClient, encrypter, privateFileWorker, privateSender),
	}
}
