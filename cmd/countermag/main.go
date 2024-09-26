package main

import (
	"context"
	"countermag/internal/logging"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func handleTest(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Test route called")

		testResponse := map[string]string{
			"message": "good test",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse)
	})
}

func main() {
	logger := logging.GetLogger("local")

	handler := http.NewServeMux()
	handler.Handle("/test", handleTest(logger))
	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	ctx := context.Background()
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
