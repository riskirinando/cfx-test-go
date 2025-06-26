package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Host      string    `json:"host"`
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/api/hello", helloHandler).Methods("GET")
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/ready", readyHandler).Methods("GET")

	// Static files (for potential frontend assets)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Go Web App</title>
		<style>
			body { font-family: Arial, sans-serif; margin: 40px; background-color: #f4f4f4; }
			.container { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
			h1 { color: #333; }
			.endpoint { background: #e8f4f8; padding: 10px; margin: 10px 0; border-radius: 4px; }
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Go Web Application</h1>
			<p>Welcome to the containerized Go web application!</p>
			<h3>Available Endpoints:</h3>
			<div class="endpoint"><strong>GET /</strong> - This page</div>
			<div class="endpoint"><strong>GET /api/hello</strong> - JSON API endpoint</div>
			<div class="endpoint"><strong>GET /health</strong> - Health check</div>
			<div class="endpoint"><strong>GET /ready</strong> - Readiness probe</div>
		</div>
	</body>
	</html>`
	fmt.Fprint(w, html)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	
	response := Response{
		Message:   "Hello from Go Web App running on Kubernetes!",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Host:      hostname,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	// Add any readiness checks here (database connectivity, etc.)
	response := HealthResponse{
		Status:    "ready",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}