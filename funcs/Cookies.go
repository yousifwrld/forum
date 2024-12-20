package forum

import (
	"context"
	"database/sql"
	"forum/db"
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

			return "", err
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
			// No cookie, no session to delete so just return
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
	_, err = db.Database.ExecContext(ctx, "DELETE FROM session WHERE sessionID = ?", sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			// if somehow there is a cookie and No rows found, no session
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

// manageSession manages the session for the user after successful login or registration
func ManageSession(w http.ResponseWriter, r *http.Request, userID int) error {
	// Delete existing session for the user
	_, err := db.Database.Exec(`DELETE FROM session WHERE userID = ?`, userID)
	if err != nil {
		log.Printf("Error deleting existing session for user ID %d: %v", userID, err)
		return err
	}

	// Set a new session cookie
	sessionID := SetCookies(w, r)

	// Insert the new session into the database
	_, err = db.Database.Exec(`INSERT INTO SESSION (SessionID, UserID) VALUES (?, ?)`, sessionID, userID)
	if err != nil {
		log.Printf("Error inserting new session for user ID %d: %v", userID, err)
		return err
	}

	return nil
}
