package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aung-arata/youtube-clone/services/user-service/internal/database"
	"github.com/aung-arata/youtube-clone/services/user-service/internal/handlers"
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

	// User routes
	userHandler := handlers.NewUserHandler(db)
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	
	// Plan routes
	planHandler := handlers.NewPlanHandler(db)
	r.HandleFunc("/plans", planHandler.GetPlans).Methods("GET")
	r.HandleFunc("/plans/{id}", planHandler.GetPlan).Methods("GET")
	r.HandleFunc("/plans", planHandler.CreatePlan).Methods("POST")
	r.HandleFunc("/plans/{id}", planHandler.UpdatePlan).Methods("PUT")
	r.HandleFunc("/plans/{id}", planHandler.DeletePlan).Methods("DELETE")
	r.HandleFunc("/users/{userId}/plan", planHandler.GetUserPlan).Methods("GET")
	r.HandleFunc("/users/{userId}/plan", planHandler.UpdateUserPlan).Methods("PUT")
	
	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"user-service"}`))
	}).Methods("GET")

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// Start server
	fmt.Printf("User Service starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
