package forum

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func SetCookies(w http.ResponseWriter, r *http.Request) string {
	sessionID := uuid.New().String()
	cookie := http.Cookie{
		Name:     "cookie",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	return sessionID
}

func GetSessionFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("cookie")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", err
		} else {

			return "", fmt.Errorf("500")
		}
	}
	return cookie.Value, nil
}

// might use this for logout, and auto session deletion
func DeleteCookiesAndSession(w http.ResponseWriter, r *http.Request) {

	// Get the session ID from the cookie
	sessionID, err := GetSessionFromCookie(r)
	if err != nil {
		if err == http.ErrNoCookie {
			// No cookie, no session
			log.Println("no cookie found, returning without deleting")
			return
		} else {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}
	}

	// Create a context with timeout for the database operation
	//which means if the database takes more than 5 seconds the operation is canceled
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Delete the session from the database
	_, err = database.ExecContext(ctx, "DELETE FROM session WHERE sessionID = ?", sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found, no session
			log.Println("no session found, returning without deleting")
			return
		}
		log.Printf("Error deleting session from database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:     "cookie",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
}
