package middleware

import (
	"net/http"
	"strings"
)

// SecurityHeadersMiddleware adds security headers to all responses
// This helps protect against common web vulnerabilities
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// X-Content-Type-Options: Prevents browsers from MIME-sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// X-Frame-Options: Prevents clickjacking attacks
		w.Header().Set("X-Frame-Options", "DENY")

		// Strict-Transport-Security: Enforces HTTPS
		// max-age=31536000 (1 year), includeSubDomains
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Content-Security-Policy: Helps prevent XSS and data injection attacks
		w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data: https:; script-src 'self'; style-src 'self' 'unsafe-inline'")

		// Referrer-Policy: Controls referrer information
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions-Policy: Controls browser features
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Cache-Control for API responses
		if strings.HasPrefix(r.URL.Path, "/api") {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		}

		next.ServeHTTP(w, r)
	})
}

// HTTPSRedirectMiddleware redirects HTTP requests to HTTPS
func HTTPSRedirectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request is already HTTPS (via X-Forwarded-Proto header or TLS)
		if r.Header.Get("X-Forwarded-Proto") == "https" || r.TLS != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Redirect to HTTPS
		httpsURL := "https://" + r.Host + r.URL.String()
		http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
	})
}
