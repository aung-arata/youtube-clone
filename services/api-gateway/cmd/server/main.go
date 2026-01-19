package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aung-arata/youtube-clone/services/api-gateway/internal/middleware"
	"github.com/gorilla/mux"
)

var (
	videoServiceURL   = getEnv("VIDEO_SERVICE_URL", "http://video-service:8081")
	userServiceURL    = getEnv("USER_SERVICE_URL", "http://user-service:8082")
	commentServiceURL = getEnv("COMMENT_SERVICE_URL", "http://comment-service:8083")
	historyServiceURL = getEnv("HISTORY_SERVICE_URL", "http://history-service:8084")
)

func main() {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Video routes - proxy to video-service
	api.PathPrefix("/videos").HandlerFunc(proxyToService(videoServiceURL, "/videos"))

	// User routes - proxy to user-service
	api.PathPrefix("/users/{id}/history").HandlerFunc(proxyToService(historyServiceURL, "/users"))
	api.PathPrefix("/users").HandlerFunc(proxyToService(userServiceURL, "/users"))

	// Comment routes - proxy to comment-service
	api.PathPrefix("/comments").HandlerFunc(proxyToService(commentServiceURL, "/comments"))
	api.PathPrefix("/videos/{videoId}/comments").HandlerFunc(proxyToService(commentServiceURL, "/videos"))

	// Health check
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"api-gateway"}`))
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
	fmt.Printf("API Gateway starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func proxyToService(serviceURL, pathPrefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the path after /api
		path := strings.TrimPrefix(r.URL.Path, "/api")

		// Build target URL
		targetURL := serviceURL + path
		if r.URL.RawQuery != "" {
			targetURL += "?" + r.URL.RawQuery
		}

		// Create new request
		proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
		if err != nil {
			http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
			return
		}

		// Copy headers
		for key, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		// Send request to target service
		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			log.Printf("Error proxying to %s: %v", targetURL, err)
			http.Error(w, "Error connecting to service", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// Set status code
		w.WriteHeader(resp.StatusCode)

		// Copy response body
		io.Copy(w, resp.Body)
	}
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

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
