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
	
	// Auth routes (public)
	authHandler := handlers.NewAuthHandler(db)
	api.HandleFunc("/auth/signup", authHandler.Signup).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST")
	
	// Protected auth route
	protectedAuth := api.PathPrefix("/auth").Subrouter()
	protectedAuth.Use(middleware.AuthMiddleware)
	protectedAuth.HandleFunc("/me", authHandler.GetCurrentUser).Methods("GET")
	
	// Video routes
	videoHandler := handlers.NewVideoHandler(db)
	api.HandleFunc("/videos", videoHandler.GetVideos).Methods("GET")
	api.HandleFunc("/videos/categories", videoHandler.GetCategories).Methods("GET")
	api.HandleFunc("/videos/trending", videoHandler.GetTrendingVideos).Methods("GET")
	api.HandleFunc("/videos/popular", videoHandler.GetPopularVideos).Methods("GET")
	api.HandleFunc("/videos/{id}", videoHandler.GetVideo).Methods("GET")
	api.HandleFunc("/videos/{id}/recommendations", videoHandler.GetRecommendations).Methods("GET")
	api.HandleFunc("/videos/{id}/analytics", videoHandler.GetVideoAnalytics).Methods("GET")
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

	// User routes
	userHandler := handlers.NewUserHandler(db)
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")

	// Watch history routes
	historyHandler := handlers.NewHistoryHandler(db)
	api.HandleFunc("/users/{userId}/history", historyHandler.AddToHistory).Methods("POST")
	api.HandleFunc("/users/{userId}/history", historyHandler.GetHistory).Methods("GET")
	
	// Plan routes
	planHandler := handlers.NewPlanHandler(db)
	api.HandleFunc("/plans", planHandler.GetPlans).Methods("GET")
	api.HandleFunc("/plans/{id}", planHandler.GetPlan).Methods("GET")
	api.HandleFunc("/plans", planHandler.CreatePlan).Methods("POST")
	api.HandleFunc("/plans/{id}", planHandler.UpdatePlan).Methods("PUT")
	api.HandleFunc("/plans/{id}", planHandler.DeletePlan).Methods("DELETE")
	api.HandleFunc("/users/{userId}/plan", planHandler.GetUserPlan).Methods("GET")
	api.HandleFunc("/users/{userId}/plan", planHandler.UpdateUserPlan).Methods("PUT")

	// Subscription routes
	subscriptionHandler := handlers.NewSubscriptionHandler(db)
	api.HandleFunc("/users/{userId}/subscriptions", subscriptionHandler.GetUserSubscriptions).Methods("GET")
	api.HandleFunc("/users/{userId}/subscriptions", subscriptionHandler.Subscribe).Methods("POST")
	api.HandleFunc("/users/{userId}/subscriptions/{channelName}", subscriptionHandler.CheckSubscription).Methods("GET")
	api.HandleFunc("/users/{userId}/subscriptions/{channelName}", subscriptionHandler.Unsubscribe).Methods("DELETE")

	// Playlist routes
	playlistHandler := handlers.NewPlaylistHandler(db)
	api.HandleFunc("/users/{userId}/playlists", playlistHandler.GetUserPlaylists).Methods("GET")
	api.HandleFunc("/users/{userId}/playlists", playlistHandler.CreatePlaylist).Methods("POST")
	api.HandleFunc("/playlists/{id}", playlistHandler.GetPlaylist).Methods("GET")
	api.HandleFunc("/playlists/{id}", playlistHandler.UpdatePlaylist).Methods("PUT")
	api.HandleFunc("/playlists/{id}", playlistHandler.DeletePlaylist).Methods("DELETE")
	api.HandleFunc("/playlists/{id}/videos", playlistHandler.AddVideoToPlaylist).Methods("POST")
	api.HandleFunc("/playlists/{id}/videos/{videoId}", playlistHandler.RemoveVideoFromPlaylist).Methods("DELETE")
	
	// Notification routes
	notificationHandler := handlers.NewNotificationHandler(db)
	api.HandleFunc("/users/{userId}/notifications", notificationHandler.GetUserNotifications).Methods("GET")
	api.HandleFunc("/users/{userId}/notifications/unread-count", notificationHandler.GetUnreadCount).Methods("GET")
	api.HandleFunc("/users/{userId}/notifications/mark-all-read", notificationHandler.MarkAllAsRead).Methods("POST")
	api.HandleFunc("/notifications", notificationHandler.CreateNotification).Methods("POST")
	api.HandleFunc("/notifications/{id}/mark-read", notificationHandler.MarkAsRead).Methods("POST")
	
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
