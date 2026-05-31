package config

import (
	"crypto/tls"
	"os"
	"testing"
)

func TestGetEnv_Default(t *testing.T) {
	os.Unsetenv("TEST_KEY_12345")
	got := getEnv("TEST_KEY_12345", "default")
	if got != "default" {
		t.Errorf("expected %q, got %q", "default", got)
	}
}

func TestGetEnv_Set(t *testing.T) {
	os.Setenv("TEST_KEY_12345", "myvalue")
	defer os.Unsetenv("TEST_KEY_12345")

	got := getEnv("TEST_KEY_12345", "default")
	if got != "myvalue" {
		t.Errorf("expected %q, got %q", "myvalue", got)
	}
}

func TestNewTLSConfig_Defaults(t *testing.T) {
	os.Unsetenv("TLS_CERT_FILE")
	os.Unsetenv("TLS_KEY_FILE")
	os.Unsetenv("TLS_ENABLED")

	cfg := NewTLSConfig()

	if cfg.CertFile != "" {
		t.Errorf("expected empty default CertFile, got %q", cfg.CertFile)
	}
	if cfg.KeyFile != "" {
		t.Errorf("expected empty default KeyFile, got %q", cfg.KeyFile)
	}
	if cfg.Enabled != false {
		t.Errorf("expected Enabled=false by default")
	}
}

func TestNewTLSConfig_EnvOverride(t *testing.T) {
	os.Setenv("TLS_CERT_FILE", "/etc/ssl/certs/server.crt")
	os.Setenv("TLS_KEY_FILE", "/etc/ssl/private/server.key")
	os.Setenv("TLS_ENABLED", "true")
	defer func() {
		os.Unsetenv("TLS_CERT_FILE")
		os.Unsetenv("TLS_KEY_FILE")
		os.Unsetenv("TLS_ENABLED")
	}()

	cfg := NewTLSConfig()

	if cfg.CertFile != "/etc/ssl/certs/server.crt" {
		t.Errorf("expected CertFile from env, got %q", cfg.CertFile)
	}
	if cfg.KeyFile != "/etc/ssl/private/server.key" {
		t.Errorf("expected KeyFile from env, got %q", cfg.KeyFile)
	}
	if !cfg.Enabled {
		t.Errorf("expected Enabled=true from env")
	}
}

func TestGetTLSConfig_MinVersion(t *testing.T) {
	tc := &TLSConfig{CertFile: "cert.pem", KeyFile: "key.pem", Enabled: true}
	tlsCfg := tc.GetTLSConfig()

	if tlsCfg == nil {
		t.Fatal("expected non-nil tls.Config")
	}
	if tlsCfg.MinVersion != tls.VersionTLS12 {
		t.Errorf("expected MinVersion TLS1.2, got %d", tlsCfg.MinVersion)
	}
}

func TestGetTLSConfig_CipherSuites(t *testing.T) {
	tc := &TLSConfig{}
	tlsCfg := tc.GetTLSConfig()

	if len(tlsCfg.CipherSuites) == 0 {
		t.Error("expected non-empty cipher suites")
	}
	// All cipher suites should be non-zero
	for _, cs := range tlsCfg.CipherSuites {
		if cs == 0 {
			t.Error("found zero cipher suite ID")
		}
	}
}
