package forum

import (
	"context"
	"database/sql"
	"log"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session cookie
		sessionID, err := GetSessionFromCookie(w, r)
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
		// Check if session exists
		row := database.QueryRow("SELECT userID FROM session WHERE sessionID = ?", sessionID)
		var userID int
		err = row.Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println(err)
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			} else {
				log.Println(err)
				ErrorPages(w, r, "500", http.StatusInternalServerError)
				return
			}
		}

		// Define a custom type for the context key
		type contextKey string
		// Set the user ID in the request context
		const userIDKey = contextKey("userID")
		ctx := r.Context()
		ctx = context.WithValue(ctx, userIDKey, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
