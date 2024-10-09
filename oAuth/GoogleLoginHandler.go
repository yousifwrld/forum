package oAuth

import (
	funcs "forum/funcs"
	"net/http"
	"net/url"
	"os"
)

var (
	googleClientID     = os.Getenv("googleClientID")
	googleClientSecret = os.Getenv("googleClientSecret")
)

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		funcs.ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}

	// Creating a new URL structure
	authURL, err := url.Parse("https://accounts.google.com/o/oauth2/v2/auth")
	if err != nil {
		funcs.ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	// Adding the required query parameters for our request
	q := authURL.Query()
	q.Add("client_id", googleClientID)
	q.Add("redirect_uri", "http://localhost:8080/googlecallback") // The callback URL which is the path they are going to send a request on for us to handle
	q.Add("response_type", "code")
	q.Add("scope", "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email")

	// Adding the query parameters to the URL
	authURL.RawQuery = q.Encode()

	// Redirecting the user to the URL
	http.Redirect(w, r, authURL.String(), http.StatusTemporaryRedirect)
}
