package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aung-arata/youtube-clone/services/api-gateway/internal/docs"
	"github.com/aung-arata/youtube-clone/services/api-gateway/internal/middleware"

	// ===== BOILERPLATE =====
	// add your telemetry package
	"github.com/aung-arata/youtube-clone/services/api-gateway/internal/telemetry"
	"github.com/gorilla/mux"

	// ===== BOILERPLATE =====
	// needed for context propagation injection in proxyToService
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var (
	videoServiceURL        = getEnv("VIDEO_SERVICE_URL", "http://video-service:8081")
	userServiceURL         = getEnv("USER_SERVICE_URL", "http://user-service:8082")
	commentServiceURL      = getEnv("COMMENT_SERVICE_URL", "http://comment-service:8083")
	historyServiceURL      = getEnv("HISTORY_SERVICE_URL", "http://history-service:8084")
	notificationServiceURL = getEnv("NOTIFICATION_SERVICE_URL", "http://notification-service:8086")
)

func main() {
	// ===== BOILERPLATE =====
	// Initialize tracer first, before anything else.
	// Non-fatal - app still works if telemetry is unavailable.
	shutdown, err := telemetry.InitTracer("api-gateway")
	if err != nil {
		log.Printf("Warning: failed to initialize tracer: %v", err)
	} else {
		// Flush remaining spans when main() exits
		defer shutdown(context.Background())
	}

	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Auth routes - proxy to user-service (public routes)
	api.PathPrefix("/auth").HandlerFunc(proxyToService(userServiceURL, "/auth"))

	// Video routes - more specific paths first
	api.PathPrefix("/upload").HandlerFunc(proxyToService(videoServiceURL, "/upload"))
	api.PathPrefix("/uploads/").HandlerFunc(proxyToService(videoServiceURL, "/uploads/"))
	api.PathPrefix("/videos/{videoId}/comments/{commentId}/replies").HandlerFunc(proxyToService(commentServiceURL, "/videos"))
	api.PathPrefix("/videos/{videoId}/comments").HandlerFunc(proxyToService(commentServiceURL, "/videos"))
	api.PathPrefix("/videos").HandlerFunc(proxyToService(videoServiceURL, "/videos"))
	api.PathPrefix("/playlists").HandlerFunc(proxyToService(videoServiceURL, "/playlists"))

	// User routes - more specific paths first to avoid catch-all /users swallowing them
	api.PathPrefix("/users/{id}/history").HandlerFunc(proxyToService(historyServiceURL, "/users"))
	api.PathPrefix("/users/{id}/subscriptions").HandlerFunc(proxyToService(userServiceURL, "/users"))
	api.PathPrefix("/users/{id}/playlists").HandlerFunc(proxyToService(videoServiceURL, "/users"))
	api.PathPrefix("/users/{id}/plan").HandlerFunc(proxyToService(userServiceURL, "/users"))

	// Notification routes - must be before generic /users catch-all
	api.PathPrefix("/users/{userId}/notifications").HandlerFunc(proxyToService(notificationServiceURL, "/users"))
	api.PathPrefix("/notifications").HandlerFunc(proxyToService(notificationServiceURL, "/notifications"))

	// User routes catch-all - proxy to user-service
	api.PathPrefix("/users").HandlerFunc(proxyToService(userServiceURL, "/users"))
	api.PathPrefix("/plans").HandlerFunc(proxyToService(userServiceURL, "/plans"))

	// Comment routes - proxy to comment-service
	api.PathPrefix("/comments").HandlerFunc(proxyToService(commentServiceURL, "/comments"))

	// API Documentation routes (Swagger/OpenAPI)
	api.HandleFunc("/docs", docs.SwaggerUIHandler).Methods("GET")
	api.HandleFunc("/docs/openapi.json", docs.OpenAPISpecHandler).Methods("GET")

	// Health check
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"api-gateway"}`))
	}).Methods("GET")

	// ===== BOILERPLATE =====
	// Middleware order matters - tracing must be first
	// so it wraps all other middleware and handlers
	r.Use(middleware.TracingMiddleware) // add this line

	// Enable CORS
	r.Use(corsMiddleware)

	// Add security headers middleware
	r.Use(middleware.SecurityHeadersMiddleware)

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

		// ===== BOILERPLATE =====
		// Inject trace context into outgoing request headers.
		// This writes traceparent header so downstream service
		// knows it belongs to the same trace.
		// r.Context() contains the active span from TracingMiddleware.
		otel.GetTextMapPropagator().Inject(
			r.Context(),
			propagation.HeaderCarrier(proxyReq.Header),
		)

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
