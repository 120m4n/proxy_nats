package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Global variable to hold version information
var versionInfo = map[string]string{
    "version": os.Getenv("APP_VERSION"),
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/health", HealthHandler).Methods("GET")
	r.HandleFunc("/api/v1/echo", EchoHandler).Methods("POST")

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add version information to the response body
	for key, value := range versionInfo {
		body[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
