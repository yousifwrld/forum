package forum

import (
	"context"
	"database/sql"
	"log"
	"net/http"
)

// custom type for the context key
type contextKey string

// Set the userID in the request context
const userIDKey = contextKey("userID")

// only for root page access so that we redirect to login if not logged in, without showing an error message
func RootAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session cookie
		sessionID, err := GetSessionFromCookie(r)
		if err != nil {
			if err == http.ErrNoCookie {
				log.Println(err)
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			} else {
				log.Println(err)
				ErrorPages(w, r, "500", http.StatusInternalServerError)
				return
			}
		}

		// Check if session exists and get userID
		userID, err := GetIDFromSession(sessionID)
		if err != nil {
			//if no session exists
			if err == sql.ErrNoRows {
				log.Println("no session found")
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			} else {
				log.Println(err)
				ErrorPages(w, r, "500", http.StatusInternalServerError)
				return
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, userIDKey, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session cookie
		sessionID, err := GetSessionFromCookie(r)
		if err != nil {
			if err == http.ErrNoCookie {
				log.Println("No cookie found:", err)
				ErrorPages(w, r, "not logged in", http.StatusUnauthorized, "templates/login.html")
				return
			} else {
				log.Println("Error getting session cookie:", err)
				ErrorPages(w, r, "500", http.StatusInternalServerError)
				return
			}
		}

		// Check if session exists and get userID
		userID, err := GetIDFromSession(sessionID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("No session found for sessionID:", sessionID)
				ErrorPages(w, r, "not logged in", http.StatusUnauthorized, "templates/login.html")
				return
			} else {
				log.Println("Error getting userID from session:", err)
				ErrorPages(w, r, "500", http.StatusInternalServerError)
				return
			}
		}

		log.Println("Authenticated userID:", userID)

		ctx := r.Context()
		ctx = context.WithValue(ctx, userIDKey, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func GetIDFromSession(sessionID string) (int, error) {
	// Check if session exists
	row := database.QueryRow("SELECT userID FROM session WHERE sessionID = ?", sessionID)
	var userID int
	err := row.Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, err
		} else {
			return 0, err
		}
	}
	return userID, nil
}