package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aung-arata/youtube-clone/services/notification-service/internal/database"
	"github.com/aung-arata/youtube-clone/services/notification-service/internal/handlers"
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

	// Notification routes
	notificationHandler := handlers.NewNotificationHandler(db)
	r.HandleFunc("/users/{userId}/notifications", notificationHandler.GetUserNotifications).Methods("GET")
	r.HandleFunc("/users/{userId}/notifications/unread-count", notificationHandler.GetUnreadCount).Methods("GET")
	r.HandleFunc("/users/{userId}/notifications/mark-all-read", notificationHandler.MarkAllAsRead).Methods("POST")
	r.HandleFunc("/notifications", notificationHandler.CreateNotification).Methods("POST")
	r.HandleFunc("/notifications/{id}/mark-read", notificationHandler.MarkAsRead).Methods("POST")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"notification-service"}`))
	}).Methods("GET")

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	// Start server
	fmt.Printf("Notification Service starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
