package forum

import (
	"net/http"

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
