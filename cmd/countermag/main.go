package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	fmt.Printf("Listening on port 8080...\n")
	server.ListenAndServe()
}