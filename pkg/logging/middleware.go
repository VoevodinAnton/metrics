package logging

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func WithLogging(next http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI

		method := r.Method

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		zap.L().Info("",
			zap.String("uri", uri),
			zap.String("method", method),
			zap.Int64("duration", int64(duration)),
		)
	}

	return http.HandlerFunc(logFn)
}
