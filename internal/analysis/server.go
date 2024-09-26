package analysis

import (
	"context"
	"countermag/internal/database"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func RunAnalysisServer(
	ctx context.Context,
	logger *slog.Logger,
	counterStore *database.Database,
	port int,
) {
	handler := http.NewServeMux()
	handler.Handle("/analysis", handleNewAnalysis(logger, counterStore))
	handler.Handle("/counts", handleGetCounts(logger, counterStore))

	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	signalCtx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	go func() {
		logger.Info(fmt.Sprintf("App server listening on %s", server.Addr), "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("error listening and serving: %s", err), "error", err)
		}

		logger.Info("Stopped serving new connections")
	}()

	<-signalCtx.Done()

	logger.Info("Shutting down application server...")

	shutdownCtx, cancel := context.WithTimeout(signalCtx, 5 * time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(fmt.Sprintf("Error shutting down application server: %s\n", err), "error", err)
	}

	logger.Info("Shutdown complete\n")
}