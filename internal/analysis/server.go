package analysis

import (
	"context"
	"countermag/internal/database"
	"countermag/internal/http/middleware"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func newHandler(
	logger *slog.Logger,
	counterStore *database.Database,
) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/analysis", handleNewAnalysis(logger, counterStore))
	mux.Handle("/counts", handleGetCounts(logger, counterStore))

	handler := middleware.AddLogging(logger, mux)

	return handler
}

func RunAnalysisServer(
	ctx context.Context,
	logger *slog.Logger,
	counterStore *database.Database,
	port int,
) {
	logger = logger.With("server", "application")

	handler := newHandler(logger, counterStore)
	server := http.Server{
		Addr:    fmt.Sprint(":" + strconv.Itoa(port)),
		Handler: handler,
	}

	go func() {
		logger.Info(fmt.Sprintf("App server listening on %s", server.Addr), "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("error listening and serving: %s", err), "error", err)
		}

		logger.Info("Stopped serving new connections")
	}()

	<-ctx.Done()

	logger.Info("Shutting down application server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(fmt.Sprintf("Error shutting down application server: %s\n", err), "error", err)
	}

	logger.Info("Shutdown complete\n")
}
