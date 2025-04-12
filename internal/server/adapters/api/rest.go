package api

import (
	"context"
	"encoding/json"
	"fmt"
	"gokeeper/internal/server/core/config"
	"gokeeper/internal/server/core/domain"
	"gokeeper/internal/server/core/logger"
	"gokeeper/pkg/auth"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-http-utils/headers"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const serverTimeout = 3

type API struct {
	srv *http.Server
}

type AuthService interface {
	Register(ctx context.Context, inUser domain.InUser) (auth.Token, error)
	Login(ctx context.Context, inUser domain.InUser) (auth.Token, error)
}

type Services interface {
	AuthService
}

type Handler struct {
	services Services
}

func NewAPI(services Services, cfg *config.Config) *API {
	h := &Handler{services}
	r := chi.NewRouter()

	r.Use(middleware.Timeout(serverTimeout * time.Second))
	r.Use(LoggingRequestMiddleware)
	r.Route("/api/user", func(r chi.Router) {
		r.Route("/register", func(r chi.Router) {
			r.Post("/", h.Register)
		})
		r.Route("/login", func(r chi.Router) {
			r.Post("/", h.Login)
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

func (h *Handler) Register(w http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Log.Debug("can not read body", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var inUser domain.InUser
	if err = json.Unmarshal(reqBody, &inUser); err != nil {
		logger.Log.Debug("can not unmarshall json", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	tokenStr, err := h.services.Register(req.Context(), inUser)
	if err != nil {
		handleException(w, err)
		return
	}
	w.Header().Set(headers.Authorization, string(tokenStr))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Log.Debug("can not read body", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var inUser domain.InUser
	if err = json.Unmarshal(reqBody, &inUser); err != nil {
		logger.Log.Debug("can not unmarshall json", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	tokenStr, err := h.services.Login(req.Context(), inUser)
	if err != nil {
		handleException(w, err)
		return
	}
	w.Header().Set(headers.Authorization, string(tokenStr))
	w.WriteHeader(http.StatusOK)

}
