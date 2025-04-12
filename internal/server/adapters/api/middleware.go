package api

import (
	"fmt"
	"gokeeper/internal/server/core/logger"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// responseData holds status and size information for responses.
type responseData struct {
	status int
	size   int
}

// loggingResponseWriter wraps an http.ResponseWriter to track response data.
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Write implements http.ResponseWriter.Write.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	if err != nil {
		return size, fmt.Errorf("failed to write response %w", err)
	}
	r.responseData.size += size
	return size, nil
}

// WriteHeader implements http.ResponseWriter.WriteHeader.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// LoggingRequestMiddleware logs incoming HTTP requests.
func LoggingRequestMiddleware(next http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		respData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   respData,
		}
		next.ServeHTTP(&lw, r)
		duration := time.Since(start)
		if respData.status == 0 {
			respData.status = 200
		}
		logger.Log.Info("got incoming http request",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.Int("status", respData.status),
			zap.Int("size", respData.size),
			zap.String("duration", duration.String()),
		)
	}
	return http.HandlerFunc(logFn)
}
