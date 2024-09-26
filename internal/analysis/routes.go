package analysis

import (
	"countermag/internal/database"
	"countermag/pkg/array"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func handleNewAnalysis(logger *slog.Logger, counterStore *database.Database) http.Handler {
	extractWords := func(r *http.Request) []string {
		defer r.Body.Close()

		var payload map[string]string
		json.NewDecoder(r.Body).Decode(&payload)

		text := payload["text"]
		words := strings.Split(text, " ")
		normalizedWords := array.MapArray(words, strings.ToLower)

		return normalizedWords
	}
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost{
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		words := extractWords(r)
		countMap := array.CountElements(words)

		logger.Info(fmt.Sprintf("Received new text containing %d new words", len(countMap)))

		for word, count := range countMap {
			counterStore.AddOccurences(word, count)
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func handleGetCounts(logger *slog.Logger, counterStore *database.Database) http.Handler {
	extractWords := func(r *http.Request) []string {
		defer r.Body.Close()

		rawWords := r.URL.Query().Get("words")
		words := strings.Split(rawWords, ",")
		return words
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		logger.Info("Received request for words")

		words := extractWords(r)

		countMap := make(map[string]int)
		for _, word := range words {
			countMap[word] = counterStore.Get(word)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(countMap)
	})
}