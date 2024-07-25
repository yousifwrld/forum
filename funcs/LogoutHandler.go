package forum

import (
	"context"
	"log"
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}

	// Get the session ID from the cookie
	sessionID, err := GetSessionFromCookie(r)
	if err != nil {
		if err.Error() == "no cookie found" {
			http.Redirect(w, r, "/login", http.StatusFound)
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
		log.Printf("Error deleting session from database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//changing the maxage to -1 to delete the cookie
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

	http.Redirect(w, r, "/login", http.StatusFound)
}
