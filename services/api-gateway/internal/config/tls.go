package config

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"
)

// TLSConfig provides TLS configuration for HTTPS servers
type TLSConfig struct {
	CertFile string
	KeyFile  string
	Enabled  bool
}

// NewTLSConfig creates a new TLS configuration from environment variables
func NewTLSConfig() *TLSConfig {
	return &TLSConfig{
		CertFile: getEnv("TLS_CERT_FILE", ""),
		KeyFile:  getEnv("TLS_KEY_FILE", ""),
		Enabled:  getEnv("TLS_ENABLED", "false") == "true",
	}
}

// GetTLSConfig returns a secure TLS configuration
func (c *TLSConfig) GetTLSConfig() *tls.Config {
	return &tls.Config{
		// Use TLS 1.2 as minimum version
		MinVersion: tls.VersionTLS12,
		// Prefer server cipher suites
		PreferServerCipherSuites: true,
		// Secure cipher suites
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		// Elliptic curves
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
	}
}

// StartServer starts either an HTTP or HTTPS server based on configuration
func (c *TLSConfig) StartServer(handler http.Handler, port string) error {
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if c.Enabled && c.CertFile != "" && c.KeyFile != "" {
		server.TLSConfig = c.GetTLSConfig()
		log.Printf("Starting HTTPS server on port %s...\n", port)
		return server.ListenAndServeTLS(c.CertFile, c.KeyFile)
	}

	log.Printf("Starting HTTP server on port %s...\n", port)
	return server.ListenAndServe()
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
