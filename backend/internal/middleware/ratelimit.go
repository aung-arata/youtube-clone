package middleware

import (
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

// RateLimitMiddleware limits requests per IP
// NOTE: This implementation uses in-memory storage and is suitable for development
// and single-instance deployments. For production with multiple instances,
// consider using Redis or a distributed cache with TTL support.
func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	// Clean up old visitors periodically
	go cleanupVisitors()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			mu.Lock()
			v, exists := visitors[ip]
			if !exists {
				visitors[ip] = &visitor{
					lastSeen: time.Now(),
					count:    1,
				}
				mu.Unlock()
				next.ServeHTTP(w, r)
				return
			}

			// Reset count if more than a minute has passed
			if time.Since(v.lastSeen) > time.Minute {
				v.count = 1
				v.lastSeen = time.Now()
				mu.Unlock()
				next.ServeHTTP(w, r)
				return
			}

			// Check if rate limit exceeded
			if v.count >= requestsPerMinute {
				mu.Unlock()
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			v.count++
			v.lastSeen = time.Now()
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}
