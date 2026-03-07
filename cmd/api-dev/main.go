// Dev-only HTTP server that wraps the Vercel handler for local testing.
// Usage: go run ./cmd/api-dev
// Listens on :8080 and serves POST /api/generate.
package main

import (
	"log"
	"net/http"

	handler "github.com/awslabs/diagram-as-code/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/generate", handler.Handler)
	log.Println("API dev server → http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
