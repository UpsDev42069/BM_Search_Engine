package weather

import (
    "testing"

    "github.com/UpsDev42069/BM_Search_Engine/backend/config"
)

func TestGetWeather(t *testing.T) {
    // Load environment variables is already handled in TestMain

    apiKey := config.APIKey
    if apiKey == "" {
        t.Fatal("API_KEY is not set in .env file")
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
