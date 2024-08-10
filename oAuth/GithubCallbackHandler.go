package oAuth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forum/db"
	funcs "forum/funcs"
	"log"
	"net/http"
)

// GithubCallbackHandler handles the GitHub OAuth callback
func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("Error: Received non-GET request for GitHub callback")
		funcs.ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}

	// Get the code from the request query param "code"
	code := r.URL.Query().Get("code")
	if code == "" {
		log.Println("Error: Missing code parameter in GitHub callback")
		funcs.ErrorPages(w, r, "400", http.StatusBadRequest)
		return
	}

	// Exchange the code for an access token
	accessToken, err := getAccessToken(code)
	if err != nil {
		log.Println("Error getting access token:", err)
		funcs.ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	// Get the user information using the access token
	email, username, githubID, err := getUserInfo(accessToken)
	if err != nil {
		log.Println("Error getting user info:", err)
		funcs.ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	// Check if user exists in the database
	exists, err := userExists("github", githubID)
	if err != nil {
		log.Println("Error checking if user exists:", err)
		funcs.ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	var userID int
	if !exists {
		// Create a new user if they don't exist
		userID, err = createUserByGithub(email, username, "github", githubID)
		if err != nil {
			log.Println("Error creating user:", err)
			funcs.ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}
	} else {
		// Get the existing user's ID
		err = db.Database.QueryRow(`SELECT UserID FROM USER WHERE oauth_userID = ?`, githubID).Scan(&userID)
		if err != nil {
			log.Println("Error getting user ID:", err)
			funcs.ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}
	}

	// Manage session
	err = funcs.ManageSession(w, r, userID)
	if err != nil {
		log.Println("Error managing session:", err)
		funcs.ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	// Redirect the user to the home page
	http.Redirect(w, r, "/", http.StatusFound)
}

// getAccessToken exchanges the code for an access token
func getAccessToken(code string) (string, error) {
	// Request structure for exchanging code for access token
	request := struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
	}{
		ClientID:     githubClientID,
		ClientSecret: githubClientSecret,
		Code:         code,
	}

	reqJson, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error marshalling request for access token exchange: %v", err)
		return "", err
	}

	// Create a new POST request to GitHub's access token endpoint
	req, err := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(reqJson),
	)
	if err != nil {
		log.Printf("Error creating request for access token exchange: %v", err)
		return "", err
	}

	// Set the necessary headers for the request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send the request and get the response
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error getting response for access token exchange: %v", err)
		return "", err
	}
	defer res.Body.Close()

	// Check for unexpected HTTP status code
	if res.StatusCode != http.StatusOK {
		log.Printf("Unexpected response status: %s", res.Status)
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response into a struct
	var response struct {
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		log.Printf("Error decoding response for access token exchange: %v", err)
		return "", err
	}

	if response.AccessToken == "" {
		log.Println("No access token in response")
		return "", fmt.Errorf("empty access token in response")
	}

	return response.AccessToken, nil
}

// getUserInfo retrieves user information using the access token
func getUserInfo(accessToken string) (string, string, int, error) {
	// Create a new GET request to GitHub's user API
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		log.Printf("Error creating request for user info: %v", err)
		return "", "", 0, err
	}

	// Set the Authorization header with the access token
	AuthHeader := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", AuthHeader)

	// Send the request and get the response
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error getting response for user info: %v", err)
		return "", "", 0, err
	}
	defer res.Body.Close()

	// Check for unexpected HTTP status code
	if res.StatusCode != http.StatusOK {
		log.Printf("Unexpected response status: %s", res.Status)
		return "", "", 0, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response into a struct
	var response struct {
		Login    string `json:"login"`
		Email    string `json:"email"`
		GithubID int    `json:"id"`
	}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		log.Printf("Error decoding response for user info: %v", err)
		return "", "", 0, err
	}

	return response.Email, response.Login, response.GithubID, nil
}

// userExists checks if a user exists by their OAuth provider and user ID
func userExists(oauthProvider string, oauthUserID int) (bool, error) {
	var exists bool
	err := db.Database.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM user WHERE oauth_provider = ? AND oauth_userID = ?)`,
		oauthProvider, oauthUserID,
	).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if user exists: %v", err)
		return false, err
	}
	return exists, nil
}

// createUserByGithub creates a new user with GitHub OAuth information
func createUserByGithub(email, username, oauthProvider string, oauthUserID int) (int, error) {
	// Start a database transaction for atomicity
	tx, err := db.Database.Begin()
	if err != nil {
		log.Printf("Error starting transaction for user creation: %v", err)
		return 0, err
	}
	defer tx.Rollback()

	// Insert the new user into the database
	result, err := tx.Exec(
		`INSERT INTO user (email, username, oauth_provider, oauth_userID) VALUES (?, ?, ?, ?)`,
		email, username, oauthProvider, oauthUserID,
	)
	if err != nil {
		log.Printf("Error inserting new user: %v", err)
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last inserted user ID: %v", err)
		return 0, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction for user creation: %v", err)
		return 0, err
	}

	return int(userID), nil
}
