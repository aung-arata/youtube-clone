package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aung-arata/youtube-clone/services/video-service/internal/database"
	"github.com/aung-arata/youtube-clone/services/video-service/internal/handlers"
	"github.com/aung-arata/youtube-clone/services/video-service/internal/middleware"
	"github.com/aung-arata/youtube-clone/services/video-service/internal/storage"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize database connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize file storage
	fileStorage, err := storage.NewFileStorage("")
	if err != nil {
		log.Fatal("Failed to initialize file storage:", err)
	}

	// Create router
	r := mux.NewRouter()

	// Serve uploaded files
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Upload routes (protected)
	uploadHandler := handlers.NewUploadHandler(db, fileStorage)
	protectedUpload := r.PathPrefix("/upload").Subrouter()
	protectedUpload.Use(middleware.AuthMiddleware)
	protectedUpload.HandleFunc("/video", uploadHandler.UploadVideo).Methods("POST")
	protectedUpload.HandleFunc("/video/delete", uploadHandler.DeleteVideo).Methods("DELETE")

	// Video routes
	videoHandler := handlers.NewVideoHandler(db)
	r.HandleFunc("/videos", videoHandler.GetVideos).Methods("GET")
	r.HandleFunc("/videos/categories", videoHandler.GetCategories).Methods("GET")
	r.HandleFunc("/videos/trending", videoHandler.GetTrendingVideos).Methods("GET")
	r.HandleFunc("/videos/popular", videoHandler.GetPopularVideos).Methods("GET")
	r.HandleFunc("/videos/{id}", videoHandler.GetVideo).Methods("GET")
	r.HandleFunc("/videos/{id}/recommendations", videoHandler.GetRecommendations).Methods("GET")
	r.HandleFunc("/videos/{id}/analytics", videoHandler.GetVideoAnalytics).Methods("GET")
	r.HandleFunc("/videos", videoHandler.CreateVideo).Methods("POST")
	r.HandleFunc("/videos/{id}/views", videoHandler.IncrementViews).Methods("POST")
	r.HandleFunc("/videos/{id}/like", videoHandler.LikeVideo).Methods("POST")
	r.HandleFunc("/videos/{id}/dislike", videoHandler.DislikeVideo).Methods("POST")

	// Playlist routes
	playlistHandler := handlers.NewPlaylistHandler(db)
	r.HandleFunc("/users/{userId}/playlists", playlistHandler.GetUserPlaylists).Methods("GET")
	r.HandleFunc("/users/{userId}/playlists", playlistHandler.CreatePlaylist).Methods("POST")
	r.HandleFunc("/playlists/{id}", playlistHandler.GetPlaylist).Methods("GET")
	r.HandleFunc("/playlists/{id}", playlistHandler.UpdatePlaylist).Methods("PUT")
	r.HandleFunc("/playlists/{id}", playlistHandler.DeletePlaylist).Methods("DELETE")
	r.HandleFunc("/playlists/{id}/videos", playlistHandler.AddVideoToPlaylist).Methods("POST")
	r.HandleFunc("/playlists/{id}/videos/{videoId}", playlistHandler.RemoveVideoFromPlaylist).Methods("DELETE")
	
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
