package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

var (
	APIKey        string
	FrontendURL   string
	SessionSecret string
	DBDriver      string // Added to specify the DB driver
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
)

// LoadEnv loads environment variables from the .env file located in the backend directory.
func LoadEnv() {
	// Determine the path to the .env file relative to this file's location.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("Unable to determine the file path.")
	}
	configDir := filepath.Dir(filename)
	envPath := filepath.Join(configDir, "../.env")

	// Convert to absolute path
	absEnvPath, err := filepath.Abs(envPath)
	if err != nil {
		log.Fatalf("Error resolving absolute path for .env file: %v", err)
	}

	// Load the .env file
	if err := godotenv.Load(absEnvPath); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Assign environment variables to package-level variables
	APIKey = os.Getenv("API_KEY")
	FrontendURL = os.Getenv("FRONTEND_URL")
	SessionSecret = os.Getenv("SESSION_SECRET")
	DBDriver = os.Getenv("DB_DRIVER") // e.g., "postgres", "sqlite3"
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")

	// Check if required environment variables are set
	missingVars := false
	if APIKey == "" {
		log.Println("Missing environment variable: API_KEY")
		missingVars = true
	}
	if FrontendURL == "" {
		log.Println("Missing environment variable: FRONTEND_URL")
		missingVars = true
	}
	if SessionSecret == "" {
		log.Println("Missing environment variable: SESSION_SECRET")
		missingVars = true
	}
	if DBDriver == "" {
		log.Println("Missing environment variable: DB_DRIVER")
		missingVars = true
	}
	if DBDriver != "sqlite3" { // Only check these variables for non-sqlite
		if DBHost == "" {
			log.Println("Missing environment variable: DB_HOST")
			missingVars = true
		}
		if DBPort == "" {
			log.Println("Missing environment variable: DB_PORT")
			missingVars = true
		}
		if DBUser == "" {
			log.Println("Missing environment variable: DB_USER")
			missingVars = true
		}
		if DBPassword == "" {
			log.Println("Missing environment variable: DB_PASSWORD")
			missingVars = true
		}
		if DBName == "" {
			log.Println("Missing environment variable: DB_NAME")
			missingVars = true
		}
	}

	if missingVars {
		log.Fatal("One or more required environment variables are not set in .env file")
	}
}