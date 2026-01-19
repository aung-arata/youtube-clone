package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aung-arata/youtube-clone/services/comment-service/internal/database"
	"github.com/aung-arata/youtube-clone/services/comment-service/internal/handlers"
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

	// Comment routes
	commentHandler := handlers.NewCommentHandler(db)
	r.HandleFunc("/videos/{videoId}/comments", commentHandler.GetComments).Methods("GET")
	r.HandleFunc("/videos/{videoId}/comments", commentHandler.CreateComment).Methods("POST")
	r.HandleFunc("/comments/{id}", commentHandler.GetComment).Methods("GET")
	r.HandleFunc("/comments/{id}", commentHandler.UpdateComment).Methods("PUT")
	r.HandleFunc("/comments/{id}", commentHandler.DeleteComment).Methods("DELETE")
	
	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"comment-service"}`))
	}).Methods("GET")

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	// Start server
	fmt.Printf("Comment Service starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
