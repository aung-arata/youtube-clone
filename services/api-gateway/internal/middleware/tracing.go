package middleware

import (
	"net/http"

	// otel is the global entry point - gives us the tracer
	"go.opentelemetry.io/otel"
	// attribute lets you attach key/value metadata to spans
	"go.opentelemetry.io/otel/attribute"
	// codes defines span status values (OK, Error)
	"go.opentelemetry.io/otel/codes"
	// propagation handles reading/writing traceparent headers
	"go.opentelemetry.io/otel/propagation"
)

func TracingMiddleware(next http.Handler) http.Handler {
	// ===== BOILERPLATE =====
	// Change "api-gateway" to your service name
	tracer := otel.Tracer("api-gateway")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ===== BOILERPLATE =====
		// Extract incoming trace context - do not change this
		ctx := otel.GetTextMapPropagator().Extract(
			r.Context(),
			propagation.HeaderCarrier(r.Header),
		)

		// ===== BOILERPLATE =====
		// Start span - span name convention is up to you but METHOD /path is standard
		spanName := r.Method + " " + r.URL.Path
		ctx, span := tracer.Start(ctx, spanName)
		defer span.End()

		// ===== PROJECT SPECIFIC =====
		// Add whatever attributes matter for your domain.
		// These are standard HTTP ones - good default for any HTTP service.
		// Add more below based on what you need to debug in your project:
		// e.g. user ID from JWT, tenant ID, video ID from path, etc.
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
			attribute.String("http.route", r.URL.Path),
			// --- add project specific attributes below this line ---
		)

		// ===== BOILERPLATE =====
		// Wrap ResponseWriter and call next - do not change this
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapped, r.WithContext(ctx))

		// ===== BOILERPLATE =====
		// Record response status - do not change this
		span.SetAttributes(
			attribute.Int("http.status_code", wrapped.statusCode),
		)

		// ===== PROJECT SPECIFIC =====
		// Error threshold - 500 is standard but you might want 400+ depending on your API
		if wrapped.statusCode >= 500 {
			span.SetStatus(codes.Error, http.StatusText(wrapped.statusCode))
		}
	})
}

// ===== BOILERPLATE =====
// responseWriter and WriteHeader - copy as-is to every project
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
