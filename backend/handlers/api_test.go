package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/UpsDev42069/BM_Search_Engine/backend/config"
    "github.com/UpsDev42069/BM_Search_Engine/backend/handlers"
    "github.com/UpsDev42069/BM_Search_Engine/backend/security"
    "github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
    // Set the environment to test
    os.Setenv("GO_ENV", "test")

    // Load environment variables
    config.LoadEnv()

	// Initialize security store with the session secret
	security.InitializeStore(config.SessionSecret)

    // Run tests
    code := m.Run()
    os.Exit(code)
}

func TestLoginHandler_Success(t *testing.T) {
    // Remove redundant godotenv.Load()
    // if err := godotenv.Load(".env"); err != nil {
    //     t.Error("Error loading .env file 1")
    // }

    // Mock database
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    username := "testuser"
    password := "password123"
    hashedPassword, _ := security.HashPassword(password)

    // Expect query for user
    mock.ExpectQuery("SELECT username, password, password_reset_required FROM users WHERE username = \\$1").
        WithArgs(username).
        WillReturnRows(sqlmock.NewRows([]string{"username", "password", "password_reset_required"}).
            AddRow(username, hashedPassword, false))

    // Prepare request
    loginReq := handlers.LoginRequest{
        Username: username,
        Password: password,
    }
    reqBody, _ := json.Marshal(loginReq)
    req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(reqBody))
    w := httptest.NewRecorder()

    // Call handler
    handler := handlers.LoginHandler(db)
    handler(w, req)

    // Assertions
    res := w.Result()
    defer res.Body.Close()
    assert.Equal(t, http.StatusOK, res.StatusCode)

    var authResp handlers.AuthResponse
    err = json.NewDecoder(res.Body).Decode(&authResp)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, authResp.StatusCode)
    assert.Equal(t, "Login successful", authResp.Message)
    assert.Equal(t, username, authResp.Username)
    assert.Equal(t, false, authResp.ResetPassword)

    // Ensure all expectations were met
    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}

func TestAPI(t *testing.T) {
    t.Skip("Not implemented")
}

func TestAPIAuth(t *testing.T) {
    t.Skip("Not implemented")
}