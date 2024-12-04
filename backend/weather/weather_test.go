package weather

import (
    "testing"
    "os"

    "github.com/UpsDev42069/BM_Search_Engine/backend/config"
)

func TestMain(m *testing.M) {
    // Set the environment to test
    os.Setenv("GO_ENV", "test")

    // Load environment variables
    config.LoadEnv()

    // Run tests
    code := m.Run()
    os.Exit(code)
}

func TestGetWeather(t *testing.T) {
    // Ensure that the API Key is loaded
    apiKey := config.APIKey
    if apiKey == "" {
        t.Fatal("API_KEY is not set in the configuration")
    }

    t.Logf("Using API_KEY: %s", apiKey) // Verifying that the API Key is loaded correctly

    // Perform the actual API call
    weatherResponse, err := GetWeather("Copenhagen", apiKey)
    if err != nil {
        t.Fatalf("Error fetching weather: %v", err)
    }

    // Logging the full response for debugging purposes
    t.Logf("Full weather response: %+v", weatherResponse)

    // Assertions to ensure we got valid data
    if weatherResponse.Name == "" {
        t.Fatalf("Expected city name, got empty string")
    }
}
