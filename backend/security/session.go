package security

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		fmt.Println("SESSION_SECRET is not set")
		os.Exit(1)
	}
	store = sessions.NewCookieStore([]byte(secret))
}

// CreateSession creates a new session with userID
func CreateSession(w http.ResponseWriter, r *http.Request, userID string) error {
	session, err := store.Get(r, "user-session")
	if err != nil {
		return err
	}

	session.Values["authenticated"] = true
	session.Values["user_id"] = userID

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	fmt.Println("Session created: ", session)

	return session.Save(r, w)
}

// GetSession retrieves the session and returns the userID
func GetSession(r *http.Request) (string, bool, error) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		return "", false, err
	}

	userID, ok := session.Values["user_id"].(string)
	authenticated, _ := session.Values["authenticated"].(bool)

	if !authenticated || !ok {
		return "", false, nil // No session found
	}

	fmt.Println("Session retrieved: ", session)

	return userID, true, nil
}

// DestroySession destroys the current session
func DestroySession(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "user-session")
	if err != nil {
		return err
	}

	session.Values["authenticated"] = false
	session.Values["user_id"] = ""

	session.Options.MaxAge = -1

	fmt.Println("Session destroyed: ", session)
	return session.Save(r, w)
}
