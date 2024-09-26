package logging

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

var TraceKey = "traceId"
var CloudTraceHeaderKey = "x-cloud-trace-context"

func getCloudOpts() *slog.HandlerOptions {
	const LevelCritical = slog.Level(12)

	opts := slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}

			if a.Key == slog.SourceKey {
				a.Key = "logging.googleapis.com/sourceLocation"
			}

			if a.Key == slog.LevelKey {
				a.Key = "severity"
				level := a.Value.Any().(slog.Level)
				if level == LevelCritical {
					a.Value = slog.StringValue("CRITICAL")
				}
			}

			if a.Key == TraceKey {
				projectId := "countermag"
				traceContext := a.Value.String()

				a.Key = "logging.googleapis.com/trace"

				traceComponents := strings.Split(traceContext, "/")

				traceId := traceComponents[0]

				value := fmt.Sprintf("projects/%s/traces/%s", projectId, traceId)

				a.Value = slog.StringValue(value)
			}

			return a
		},
	}

	return &opts
}

func GetLogger(environment string) *slog.Logger {
	if environment == "local" {
		logger := slog.New(newHandler(nil))

		return logger
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, getCloudOpts()))
}

var Logger = GetLogger("local")
