package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"
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
	srv := &http.Server{Addr: ":" + port, Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("Shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}
