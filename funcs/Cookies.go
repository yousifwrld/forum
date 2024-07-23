package forum

import (
	"fmt"
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

func GetSessionFromCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("cookie")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("no cookie found")
		} else {

			return "", fmt.Errorf("500")
		}
	}
	return cookie.Value, nil
}
