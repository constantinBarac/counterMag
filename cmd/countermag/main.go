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
	
	closeCtx, cancel := context.WithTimeout(ctx, 6 * time.Second)
	
	analysis.RunAnalysisServer(
		signalCtx,
		logger,
		counterStore,
		8080,
	)

	counterStore.Close(closeCtx)
	defer cancel()
}
