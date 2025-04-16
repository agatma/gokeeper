package api

import (
	"context"
	"fmt"
	"gokeeper/internal/server/core/config"
	"gokeeper/internal/server/core/domain"
	"gokeeper/pkg/auth"
	"gokeeper/pkg/logger"
	"gokeeper/pkg/middlewares"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const serverTimeout = 3

type API struct {
	srv *http.Server
}

type AuthService interface {
	Register(ctx context.Context, inUser domain.InUserRequest) (auth.Token, error)
	Login(ctx context.Context, inUser domain.InUserRequest) (auth.Token, error)
}

type PrivateService interface {
	Save(ctx context.Context, pd *domain.Data, userID uuid.UUID) error
	GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain.Data, error)
	Delete(ctx context.Context, pd *domain.DeleteRequest, userID uuid.UUID) error
	GetAll(ctx context.Context, req *domain.GetAllRequest, userID uuid.UUID) ([]domain.Data, error)
}

type Services interface {
	AuthService
	PrivateService
}

type Handler struct {
	services Services
}

func NewAPI(services Services, cfg *config.Config, auth *auth.Authenticator) *API {
	h := &Handler{services}
	r := chi.NewRouter()

	r.Use(middleware.Timeout(serverTimeout * time.Second))
	r.Use(middlewares.LoggingRequestMiddleware)
	r.Route("/api/user", func(r chi.Router) {
		r.Route("/register", func(r chi.Router) {
			r.Post("/", h.Register)
		})
		r.Route("/login", func(r chi.Router) {
			r.Post("/", h.Login)
		})
	})
	r.Route("/api/private", func(r chi.Router) {
		r.Use(middlewares.AuthenticateMiddleware(auth))
		r.Group(func(r chi.Router) {
			r.Post("/", h.Save)
			r.Delete("/", h.Delete)
		})
		r.Group(func(r chi.Router) {
			r.Get("/{id:^[a-zA-Z0-9-_]+}", h.Get)
			r.Group(func(r chi.Router) {
				r.Get("/", h.GetAll)
			})
		})
	})
	return &API{
		srv: &http.Server{
			Addr:    cfg.Address,
			Handler: r,
		},
	}
}

// Run starts the HTTP server.
func (a *API) Run() error {
	sigint := make(chan os.Signal, 1)
	// Graceful shutdown of server
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigint
		if err := a.srv.Shutdown(context.Background()); err != nil {
			logger.Log.Info("server shutdown gracefully: ", zap.Error(err))
		}
	}()
	if err := a.srv.ListenAndServe(); err != nil {
		logger.Log.Error("error occurred during running server: ", zap.Error(err))
		return fmt.Errorf("failed run server: %w", err)
	}
	return nil
}
