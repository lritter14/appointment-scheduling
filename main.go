package main

import (
	"appointment-scheduling/data"
	"appointment-scheduling/routes"
	"log"
	"net/http"
)

func main() {
	// Initialize routes
	router := routes.InitializeRoutes()

	// Initialize database
	db, err := data.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
