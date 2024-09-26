package main

import (
	"context"
	"countermag/internal/analysis"
	"countermag/internal/database"
	"countermag/internal/logging"
	"os"
	"os/signal"
	"time"
)


func main() {
	ctx := context.Background()
	signalCtx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := logging.GetLogger("local")

	snapshotPath := "counter.txt"
	counterStore := database.NewDatabase(signalCtx, logger, &database.FileSnapshotPersister{Path: snapshotPath})
	
	
	analysis.RunAnalysisServer(
		signalCtx,
		logger,
		counterStore,
		8080,
	)
	
	closeCtx, cancel := context.WithTimeout(ctx, 5 * time.Second)
	counterStore.Close(closeCtx)
	defer cancel()
}
