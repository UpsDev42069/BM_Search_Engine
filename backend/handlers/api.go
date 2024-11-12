package handlers

// Imports for the handlers
import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"github.com/UpsDev42069/BM_Search_Engine/backend/db"
	"github.com/UpsDev42069/BM_Search_Engine/backend/security"
	"github.com/UpsDev42069/BM_Search_Engine/backend/weather"
	"github.com/joho/godotenv"
)

type AuthResponse struct {
	StatusCode    int    `json:"statusCode"`
	Message       string `json:"message"`
	Username      string `json:"username,omitempty"`      // Add the Username field
	ResetPassword bool   `json:"resetPassword,omitempty"` // Add the ResetPassword field
}

type SearchResponse struct {
	Data []map[string]interface{} `json:"data"`
}

type StandardResponse struct {
	Data map[string]interface{} `json:"data"`
}

type RequestValidationError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type HTTPValidationError struct {
	Detail []ValidationError `json:"detail"`
}

type ValidationError struct {
	Loc  []interface{} `json:"loc"`
	Msg  string        `json:"msg"`
	Type string        `json:"type"`
}

// RootGet handles the root GET request
func RootGet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

// @Summary Search
// @Description Search for content
// @Tags search
// @Accept  json
// @Produce  json
// @Param q query string true "Query"
// @Param language query string false "Language code (e.g., 'en')"
// @Success 200 {object} SearchResponse
// @Failure 422 {object} RequestValidationError
// @Router /api/search [get]
func SearchHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(RequestValidationError{
				StatusCode: 422,
				Message:    "Query parameter is required",
			})
			return
		}

		language := r.URL.Query().Get("language")
		if language == "" {
			language = "en"
		}

		var searchResults []map[string]interface{}
		// Changed to use FTS5 in sqlite3
		query := `
    SELECT pages.title, pages.url, pages.language, pages.last_updated, pages.content
    FROM pages_fts
    JOIN pages ON pages_fts.title = pages.title
    WHERE pages_fts MATCH ?
		`
		args := []interface{}{q}
		results, err := db.QueryDB(database, query, args...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		searchResults = results

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResponse{
			Data: searchResults,
		})
	}
}

// RootPost handles the root POST request
func RootPost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

/////////////////////////////////////////
// Funtions for registering a new user //
/////////////////////////////////////////

// User represents the user registration data
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// @Summary Register
// @Description Register a new user
// @Tags register
// @Accept  application/json
// @Produce  json
// @Param username formData string true "Username"
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Param password2 formData string true "Password2"
// @Success 200 {object} AuthResponse
// @Failure 422 {object} HTTPValidationError
// @Router /api/register [post]
func RegisterHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var user RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if user.Username == "" || user.Password == "" || user.Email == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(HTTPValidationError{
				Detail: []ValidationError{
					{
						Loc:  []interface{}{"body", "username"},
						Msg:  "field required",
						Type: "value_error.missing",
					},
					{
						Loc:  []interface{}{"body", "password"},
						Msg:  "field required",
						Type: "value_error.missing",
					},
					{
						Loc:  []interface{}{"body", "email"},
						Msg:  "field required",
						Type: "value_error.missing",
					},
				},
			})
			return
		}

		// Save the user to the database
		insertUserSQL := `INSERT INTO users (username, password, email) VALUES (?, ?, ?)`

		// Hash password before saving it to the database
		hashedPassword, err := security.HashPassword(user.Password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		_, err = database.Exec(insertUserSQL, user.Username, hashedPassword, user.Email)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AuthResponse{StatusCode: http.StatusOK, Message: "User registered successfully"})
	}
}

// Error logic if search query is empty
func SearchHandlerLucas(w http.ResponseWriter, r *http.Request) {
	// function body
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Searching for " + query))
}

/////////////////////////////////////////
//	 	Login logic for users        //
/////////////////////////////////////////

// LoginRequest represents the user login data
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary Login
// @Description Login a user
// @Tags login
// @Accept  application/json
// @Produce  json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} AuthResponse
// @Failure 422 {object} HTTPValidationError
// @Router /api/login [post]
func LoginHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var user LoginRequest
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if user.Username == "" || user.Password == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(HTTPValidationError{
				Detail: []ValidationError{
					{
						Loc:  []interface{}{"body", "username"},
						Msg:  "field required",
						Type: "value_error.missing",
					},
					{
						Loc:  []interface{}{"body", "password"},
						Msg:  "field required",
						Type: "value_error.missing",
					},
				},
			})
			return
		}

		// Query the database for the user
		var dbUser LoginRequest
		var resetPassword bool
		err = database.QueryRow("SELECT username, password, password_reset_required FROM users WHERE username = ?", user.Username).Scan(&dbUser.Username, &dbUser.Password, &resetPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Validate the password
		if !security.CheckPasswordHash(dbUser.Password, user.Password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		//Create a new session
		err = security.CreateSession(w, r, dbUser.Username)
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AuthResponse{
			StatusCode:    http.StatusOK,
			Message:       "Login successful",
			Username:      dbUser.Username,
			ResetPassword: resetPassword,
		})
	}
}

// @Summary Weather
// @Description Get weather information
// @Tags weather
// @Produce  json
// @Success 200 {object} StandardResponse
// @Router /api/weather [get]
func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	if err := godotenv.Load(); err != nil {
		http.Error(w, "Error loading .env file", http.StatusInternalServerError)
		return
	}
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		http.Error(w, "API_KEY is not set in .env file", http.StatusInternalServerError)
		return
	}

	// Fetch the weather data for the city
	weatherResponse, err := weather.GetWeather("Copenhagen", apiKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Set response header and send JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weatherResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// @Summary Logout
// @Description Logout a user
// @Tags logout
// @Produce  json
// @Success 200 {object} AuthResponse
// @Router /api/logout [get]
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := security.DestroySession(w, r)
	if err != nil {
		http.Error(w, "Failed to destroy session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{StatusCode: http.StatusOK, Message: "Logout successful"})
}

// @Summary Check Login
// @Description Check if the user is logged in
// @Tags auth
// @Produce  json
// @Success 200 {object} AuthResponse
// @Router /api/check-login [get]
func CheckLoginHandler(w http.ResponseWriter, r *http.Request) {
	userID, authenticated, err := security.GetSession(r)
	if err != nil {
		http.Error(w, "Error checking session", http.StatusInternalServerError)
		return
	}

	if !authenticated || userID == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AuthResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Not logged in",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		StatusCode: http.StatusOK,
		Message:    "Logged in",
		Username:   userID,
	})
}

type ResetPasswordRequest struct {
	Username    string `json:"username"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// @Summary Reset password
func ResetPasswordHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var user ResetPasswordRequest
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var hashedPassword string
		err = database.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Check if the old password is correct

		if !security.CheckPasswordHash(hashedPassword, user.OldPassword) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Hash the new password and update the user in the database
		newHashedPassword, err := security.HashPassword(user.NewPassword)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// Update the password in the database
		_, err = database.Exec("UPDATE users SET password = ?, password_reset_required = ? WHERE username = ?", newHashedPassword, false, user.Username)
		if err != nil {
			http.Error(w, "Failed to update password", http.StatusInternalServerError)
			return
		}

		//Return success response

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AuthResponse{
			StatusCode: http.StatusOK,
			Message:    "Password reset successful",
			Username:   user.Username,
		})

	}
}
