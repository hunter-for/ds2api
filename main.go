package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	defaultPort    = 8080
	defaultHost    = "0.0.0.0"
	appName        = "ds2api"
	appVersion     = "dev"
	defaultDSPort  = 2302 // DayZ default game port
)

// Config holds the application configuration loaded from environment variables.
type Config struct {
	Host        string
	Port        int
	DSAddress   string
	DSPort      int
	APIKey      string
	Debug       bool
}

func main() {
	// Load .env file if present (ignored in production where env vars are set directly)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if cfg.Debug {
		log.Printf("%s %s starting in debug mode", appName, appVersion)
		log.Printf("Listening on %s:%d", cfg.Host, cfg.Port)
		log.Printf("Downstream service: %s:%d", cfg.DSAddress, cfg.DSPort)
	}

	router := setupRouter(cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("%s listening on %s", appName, addr)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

// loadConfig reads configuration from environment variables with sensible defaults.
func loadConfig() (*Config, error) {
	cfg := &Config{
		Host:      getEnv("API_HOST", defaultHost),
		Port:      defaultPort,
		DSAddress: getEnv("DS_ADDRESS", "localhost"),
		DSPort:    defaultDSPort,
		APIKey:    getEnv("API_KEY", ""),
		Debug:     getEnvBool("DEBUG", false),
	}

	if portStr := os.Getenv("API_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid API_PORT value %q: %w", portStr, err)
		}
		cfg.Port = port
	}

	if portStr := os.Getenv("DS_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid DS_PORT value %q: %w", portStr, err)
		}
		cfg.DSPort = port
	}

	if cfg.DSAddress == "" {
		return nil, fmt.Errorf("DS_ADDRESS must be set")
	}

	return cfg, nil
}

// getEnv returns the value of the environment variable named by key,
// or fallback if the variable is not set or empty.
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// getEnvBool returns the boolean value of the environment variable named by key,
// or fallback if the variable is not set or cannot be parsed.
func getEnvBool(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return b
}
