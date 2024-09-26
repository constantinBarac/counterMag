package main

import (
	"context"
	"countermag/internal/analysis"
	"countermag/internal/database"
	"countermag/internal/logging"
)


func main() {
	ctx := context.Background()
	logger := logging.GetLogger("local")

	counterStore := database.NewDatabase(ctx, logger, &database.FileSnapshotPersister{})

	analysis.RunAnalysisServer(
		ctx,
		logger,
		counterStore,
		8080,
	)
}
