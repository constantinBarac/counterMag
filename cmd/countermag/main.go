package main

import (
	"context"
	"countermag/internal/analysis"
	"countermag/internal/database"
	"countermag/internal/logging"
)


func main() {
	logger := logging.GetLogger("local")

	counterStore := database.NewDatabase(logger)

	analysis.RunAnalysisServer(
		context.Background(),
		logger,
		counterStore,
		8080,
	)
}
