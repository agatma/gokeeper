package api

import (
	"encoding/json"
	"gokeeper/pkg/domain"
	"gokeeper/pkg/logger"
	"io"
	"net/http"

	"github.com/go-http-utils/headers"
	"go.uber.org/zap"
)

func (h *Handler) Register(w http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Log.Debug("can not read body", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var inUser domain.InUserRequest
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

	var inUser domain.InUserRequest
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
