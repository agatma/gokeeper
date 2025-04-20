package app

import (
	"fmt"
	"gokeeper/internal/server/adapters/api"
	"gokeeper/internal/server/adapters/storage"
	"gokeeper/internal/server/core/config"
	"gokeeper/internal/server/core/service"
	"gokeeper/pkg/auth"
	"gokeeper/pkg/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type Server struct {
	cfg      *config.Config
	api      *api.API
	services *service.Services
}

func NewServer() (*Server, error) {
	cfg := config.NewConfig()
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return nil, fmt.Errorf("can't load logger: %w", err)
	}
	newStorage, err := storage.NewStorage(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}
	authenticator := auth.NewAuthJWT(cfg.JWTSecretKey, cfg.TokenExp)
	services := service.NewServices(newStorage, *authenticator)
	return &Server{
		cfg: cfg,
		api: api.NewAPI(services, cfg, authenticator),
	}, nil
}

func (s *Server) Run() {
	if err := s.api.Run(); err != nil {
		logger.Log.Error("error while running server", zap.Error(err))
		return
	}
}
