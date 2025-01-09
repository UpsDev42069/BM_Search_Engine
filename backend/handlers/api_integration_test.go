// backend/handlers/api_integration_test.go

package handlers_test

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

    "github.com/UpsDev42069/BM_Search_Engine/backend/config"
    "github.com/UpsDev42069/BM_Search_Engine/backend/db"
    "github.com/UpsDev42069/BM_Search_Engine/backend/handlers"
    "github.com/UpsDev42069/BM_Search_Engine/backend/security"
    "github.com/stretchr/testify/assert"
)

// WeatherResponse represents the structure of the weather API response.
// Define this according to the actual response you expect.
type WeatherResponse struct {
    Name string `json:"name"`
    // Add other fields as necessary
}

type AuthResponse struct {
    StatusCode    int    `json:"statusCode"`
    Message       string `json:"message"`
    Username      string `json:"username,omitempty"`
    ResetPassword bool   `json:"resetPassword,omitempty"`
}

func TestWeatherHandler_Integration(t *testing.T) {
    // Set the environment to test
    os.Setenv("GO_ENV", "test")

    // Load environment variables
    config.LoadEnv()

    // Initialize security store with the session secret
    security.InitializeStore(config.SessionSecret)

    // Connect to the real database
    database, err := db.ConnectDB(false) // Pass 'true' if using Dockerized DB
    assert.NoError(t, err, "Database connection should not fail")
    defer database.Close()

    // Run migrations if necessary
    db.RunMigrations(database)
    if err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    // Initialize the handler
    handler := http.HandlerFunc(handlers.WeatherHandler)

    // Create a new HTTP request
    req, err := http.NewRequest("GET", "/api/weather", nil)
    assert.NoError(t, err, "Creating the request should not fail")

    // Use httptest to record the response
    rr := httptest.NewRecorder()

    // Serve the HTTP request
    handler.ServeHTTP(rr, req)

    // Check the status code
    assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")

    // Decode the response body
    var weatherResp WeatherResponse
    err = json.NewDecoder(rr.Body).Decode(&weatherResp)
    assert.NoError(t, err, "Decoding response should not fail")

    // Assertions: Verify that the response contains expected data
    assert.NotEmpty(t, weatherResp.Name, "Expected city name in response")
    // Add more assertions based on the WeatherResponse structure
}

func TestLoginHandler_Integration(t *testing.T) {
    // Set the environment to test
    os.Setenv("GO_ENV", "test")

    // Load environment variables
    config.LoadEnv()

    // Initialize security store with the session secret
    security.InitializeStore(config.SessionSecret)

    // Connect to the real database
    database, err := db.ConnectDB(false) // Pass 'true' if using Dockerized DB
    assert.NoError(t, err, "Database connection should not fail")
    defer database.Close()

    // Run migrations if necessary
    db.RunMigrations(database)
    if err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    // Ensure the test user exists
    testUsername := "integration_test_user"
    testPassword := "securepassword"
    hashedPassword, err := security.HashPassword(testPassword)
    assert.NoError(t, err, "Password hashing should not fail")

    // Insert the test user into the database
    _, err = database.Exec(`
        INSERT INTO users (username, password, password_reset_required)
        VALUES ($1, $2, $3)
        ON CONFLICT (username) DO NOTHING
    `, testUsername, hashedPassword, false)
    assert.NoError(t, err, "Inserting test user should not fail")

    // Prepare the login request payload
    loginReq := handlers.LoginRequest{
        Username: testUsername,
        Password: testPassword,
    }
    reqBody, err := json.Marshal(loginReq)
    assert.NoError(t, err, "Marshaling login request should not fail")

    // Create a new HTTP request
    req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(reqBody))
    assert.NoError(t, err, "Creating the login request should not fail")

    // Set the appropriate headers
    req.Header.Set("Content-Type", "application/json")

    // Use httptest to record the response
    rr := httptest.NewRecorder()

    // Call the handler with the real database
    handler := handlers.LoginHandler(database)
    handler.ServeHTTP(rr, req)

    // Check the status code
    assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")

    // Decode the response body
    var authResp AuthResponse
    err = json.NewDecoder(rr.Body).Decode(&authResp)
    assert.NoError(t, err, "Decoding response should not fail")

    // Assertions: Verify authentication was successful
    assert.Equal(t, http.StatusOK, authResp.StatusCode, "Expected status code 200 in response")
    assert.Equal(t, "Login successful", authResp.Message, "Expected successful login message")
    assert.Equal(t, testUsername, authResp.Username, "Expected correct username in response")
    assert.False(t, authResp.ResetPassword, "Expected reset password to be false")

    // Cleanup: Remove the test user from the database
    _, err = database.Exec("DELETE FROM users WHERE username = $1", testUsername)
    assert.NoError(t, err, "Cleaning up test user should not fail")
}

func TestRegisterHandler_Integration(t *testing.T) {
    // Set the environment to test
    os.Setenv("GO_ENV", "test")
	
    // Load environment variables
    config.LoadEnv()

    // Initialize security store with the session secret
    security.InitializeStore(config.SessionSecret)

    // Connect to the real database
    database, err := db.ConnectDB(false) // Pass 'true' if using Dockerized DB
    assert.NoError(t, err, "Database connection should not fail")
    defer database.Close()

    // Run migrations if necessary
    db.RunMigrations(database)
    if err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    // Prepare the registration request payload
    registerReq := handlers.RegisterRequest{
        Username: "integration_test_user_register",
        Password: "newsecurepassword",
    }
    reqBody, err := json.Marshal(registerReq)
    assert.NoError(t, err, "Marshaling registration request should not fail")

    // Create a new HTTP request
    req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(reqBody))
    assert.NoError(t, err, "Creating the registration request should not fail")

    // Set the appropriate headers
    req.Header.Set("Content-Type", "application/json")

    // Use httptest to record the response
    rr := httptest.NewRecorder()

    // Call the handler
    handler := handlers.RegisterHandler(database)
    handler.ServeHTTP(rr, req)

    // Check the status code
    assert.Equal(t, http.StatusCreated, rr.Code, "Expected status code 201")

    // Decode the response body
    var resp handlers.StandardResponse
    err = json.NewDecoder(rr.Body).Decode(&resp)
    assert.NoError(t, err, "Decoding response should not fail")
    assert.Equal(t, "User registered successfully", resp.Message)
    // Verify that the user exists in the database
    var exists bool
    err = database.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", registerReq.Username).Scan(&exists)
    assert.NoError(t, err, "Querying user existence should not fail")
    assert.True(t, exists, "User should exist in the database")

    // Cleanup: Remove the test user from the database
    _, err = database.Exec("DELETE FROM users WHERE username = $1", registerReq.Username)
    assert.NoError(t, err, "Cleaning up test user should not fail")
}