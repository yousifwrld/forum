package oAuth

import (
	funcs "forum/funcs"
	"net/http"
	"net/url"
)

var (
	githubClientID     = "Ov23liJQhTGw0Kzh4PRN"
	githubClientSecret = "d350cb1e955cce4ddd644c42655d50261db09dee"
)

func GithubloginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		funcs.ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}

	//creating a new url structure
	url, err := url.Parse("https://github.com/login/oauth/authorize")
	if err != nil {
		funcs.ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	//adding the required query parameters for our request
	q := url.Query()
	//the clientID we got from github
	q.Add("client_id", githubClientID)
	//the callback url which is the path they are going to send a request on for us to handle
	q.Add("redirect_uri", "http://localhost:8080/githubcallback")

	//adding the query parameters to the url
	url.RawQuery = q.Encode()
	//redirecting the user to the url
	http.Redirect(w, r, url.String(), http.StatusTemporaryRedirect)
}
