package api

import (
	"errors"
	"gokeeper/pkg/domain"
	"gokeeper/pkg/logger"
	"net/http"

	"go.uber.org/zap"
)

func handleException(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrUserConflict):
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
	case errors.Is(err, domain.ErrUserAuthentication):
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	case errors.Is(err, domain.ErrPrivateDataConflict):
		w.WriteHeader(http.StatusConflict)
	case errors.Is(err, domain.ErrPrivateDataBadFormat):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, domain.ErrPrivateDataNotFound):
		w.WriteHeader(http.StatusNotFound)
	default:
		logger.Log.Error("Internal server error", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
