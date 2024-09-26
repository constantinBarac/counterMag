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

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	go func() {
		fmt.Printf("App server listening on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("error listening and serving: %s\n", err)
		}

		fmt.Printf("Stopped serving new connections\n")
	}()

	<-ctx.Done()

	fmt.Printf("Shutting down application server...\n")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Error shutting down application server: %s\n", err)
	}

	fmt.Printf("Shutdown complete\n")
}