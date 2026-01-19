package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aung-arata/youtube-clone/services/history-service/internal/database"
	"github.com/aung-arata/youtube-clone/services/history-service/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize database connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create router
	r := mux.NewRouter()

	// History routes
	historyHandler := handlers.NewHistoryHandler(db)
	r.HandleFunc("/users/{userId}/history", historyHandler.AddToHistory).Methods("POST")
	r.HandleFunc("/users/{userId}/history", historyHandler.GetHistory).Methods("GET")
	
	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"history-service"}`))
	}).Methods("GET")

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	// Start server
	fmt.Printf("History Service starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
