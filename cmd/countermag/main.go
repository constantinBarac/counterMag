package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func handleTest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Test route called")

	testResponse := map[string]string{
		"message": "good test",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(testResponse)
}

func main() {
	handler := http.NewServeMux()
	handler.HandleFunc("/test", handleTest)
	server := http.Server{
		Addr: ":8080",
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