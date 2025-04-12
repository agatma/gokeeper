package api

import (
	"errors"
	"gokeeper/internal/server/core/domain"
	"gokeeper/internal/server/core/logger"
	"net/http"

	"go.uber.org/zap"
)

func handleException(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrUserConflict):
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
	case errors.Is(err, domain.ErrUserAuthentication):
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	default:
		logger.Log.Error("Internal server error", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
