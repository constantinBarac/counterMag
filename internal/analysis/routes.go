package analysis

import (
	"countermag/internal/database"
	"countermag/pkg/array"
	"countermag/pkg/text"
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

		payloadText, hasText := payload["text"]

		if !hasText {
			return nil
		}

		strippedAndNormalizedText := strings.ToLower(text.Strip(payloadText))
		words := strings.Split(strippedAndNormalizedText, " ")

		return words
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		words := extractWords(r)

		if words == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "Text parameter can not be empty"})
			return
		}

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

		nonEmptyWords := make([]string, 0, len(words))
		for _, word := range words {
			if word != "" {
				nonEmptyWords = append(nonEmptyWords, word)
			}
		}

		return nonEmptyWords
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		logger.Info("Received request for words")

		words := extractWords(r)
		fmt.Printf("----> %s | %d", words, len(words))
		if len(words) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "You must specify at least one word"})
			return
		}

		countMap := make(map[string]int)
		for _, word := range words {
			countMap[word] = counterStore.Get(word)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(countMap)
	})
}
