package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aung-arata/youtube-clone/backend/internal/database"
	"github.com/aung-arata/youtube-clone/backend/internal/handlers"
	"github.com/aung-arata/youtube-clone/backend/internal/middleware"
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

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	
	// Video routes
	videoHandler := handlers.NewVideoHandler(db)
	api.HandleFunc("/videos", videoHandler.GetVideos).Methods("GET")
	api.HandleFunc("/videos/{id}", videoHandler.GetVideo).Methods("GET")
	api.HandleFunc("/videos", videoHandler.CreateVideo).Methods("POST")
	api.HandleFunc("/videos/{id}/views", videoHandler.IncrementViews).Methods("POST")
	api.HandleFunc("/videos/{id}/like", videoHandler.LikeVideo).Methods("POST")
	api.HandleFunc("/videos/{id}/dislike", videoHandler.DislikeVideo).Methods("POST")
	
	// Comment routes
	commentHandler := handlers.NewCommentHandler(db)
	api.HandleFunc("/videos/{videoId}/comments", commentHandler.GetComments).Methods("GET")
	api.HandleFunc("/videos/{videoId}/comments", commentHandler.CreateComment).Methods("POST")
	api.HandleFunc("/comments/{id}", commentHandler.GetComment).Methods("GET")
	api.HandleFunc("/comments/{id}", commentHandler.UpdateComment).Methods("PUT")
	api.HandleFunc("/comments/{id}", commentHandler.DeleteComment).Methods("DELETE")
	
	// Health check
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Enable CORS
	r.Use(corsMiddleware)
	
	// Add logging middleware
	r.Use(middleware.LoggingMiddleware)
	
	// Add rate limiting middleware (100 requests per minute)
	r.Use(middleware.RateLimitMiddleware(100))

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
