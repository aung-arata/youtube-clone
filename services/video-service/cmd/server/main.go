package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aung-arata/youtube-clone/services/video-service/internal/database"
	"github.com/aung-arata/youtube-clone/services/video-service/internal/handlers"
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

	// Video routes
	videoHandler := handlers.NewVideoHandler(db)
	r.HandleFunc("/videos", videoHandler.GetVideos).Methods("GET")
	r.HandleFunc("/videos/categories", videoHandler.GetCategories).Methods("GET")
	r.HandleFunc("/videos/{id}", videoHandler.GetVideo).Methods("GET")
	r.HandleFunc("/videos", videoHandler.CreateVideo).Methods("POST")
	r.HandleFunc("/videos/{id}/views", videoHandler.IncrementViews).Methods("POST")
	r.HandleFunc("/videos/{id}/like", videoHandler.LikeVideo).Methods("POST")
	r.HandleFunc("/videos/{id}/dislike", videoHandler.DislikeVideo).Methods("POST")
	
	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"video-service"}`))
	}).Methods("GET")

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Start server
	fmt.Printf("Video Service starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
