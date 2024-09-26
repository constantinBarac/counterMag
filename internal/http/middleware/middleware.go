package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func AddLogging(logger *slog.Logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)

		duration := time.Since(start).Milliseconds()
		logger.Info(
			fmt.Sprintf("%s %s took %dms",
				r.Method, r.URL.Path, duration,
			),
			"method", r.Method,
			"path", r.URL.Path,
			"duration", duration,
		)
	})
}
