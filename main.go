package main

import (
	// "go-api-nutrishe/controllers/nabila"
	"log"
	"net/http"
	"nutrishe/controllers/nabila"
	"nutrishe/models"
)

func main() {
	// Setup database
	err := models.Setup()
	if err != nil {
		log.Fatalf("Failed to set up database: %v", err)
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Define routes
	mux.HandleFunc("/empowher", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/checkin", nabila.Checkin)

	// Start the HTTP server
	port := ":8081"
	log.Printf("Starting server on port %s", port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
