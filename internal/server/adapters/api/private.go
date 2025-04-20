package api

import (
	"encoding/json"
	"gokeeper/pkg/domain"
	"gokeeper/pkg/logger"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-http-utils/headers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *Handler) Save(w http.ResponseWriter, req *http.Request) {
	userID, err := uuid.Parse(req.Header.Get("X-User-ID"))
	if err != nil {
		logger.Log.Error("failed to parse X-User-ID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Log.Error("failed to read request body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var privateData domain.Data
	if err = json.Unmarshal(reqBody, &privateData); err != nil {
		logger.Log.Error("failed to unmarshal private data", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.services.Save(req.Context(), &privateData, userID); err != nil {
		handleException(w, err)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Delete(w http.ResponseWriter, req *http.Request) {
	userID, err := uuid.Parse(req.Header.Get("X-User-ID"))
	if err != nil {
		logger.Log.Error("failed to parse X-User-ID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var privateDeleteRequest domain.DeleteRequest
	if err = json.Unmarshal(reqBody, &privateDeleteRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.services.Delete(req.Context(), &privateDeleteRequest, userID); err != nil {
		handleException(w, err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Get(w http.ResponseWriter, req *http.Request) {
	userID, err := uuid.Parse(req.Header.Get("X-User-ID"))
	if err != nil {
		logger.Log.Error("failed to parse X-User-ID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dataId := chi.URLParam(req, "id")
	if dataId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	privateData, err := h.services.GetByID(req.Context(), dataId, userID)
	if err != nil {
		handleException(w, err)
		return
	}

	resp, err := json.Marshal(privateData)
	if err != nil {
		logger.Log.Error("failed to parse json", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(headers.ContentType, "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) GetAll(w http.ResponseWriter, req *http.Request) {
	userID, err := uuid.Parse(req.Header.Get("X-User-ID"))
	if err != nil {
		logger.Log.Error("failed to parse X-User-ID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	limit := req.URL.Query().Get("limit")
	if limit == "" {
		limit = "10"
	}
	offset := req.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}

	var GetAllRequest domain.GetAllRequest

	if GetAllRequest.Limit, err = strconv.ParseUint(limit, 0, 64); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if GetAllRequest.Offset, err = strconv.ParseUint(offset, 0, 64); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	privateData, err := h.services.GetAll(req.Context(), &GetAllRequest, userID)
	if err != nil {
		logger.Log.Error("GetAll: internal error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(privateData)
	if err != nil {
		logger.Log.Error("GetAll: internal error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(headers.ContentType, "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
