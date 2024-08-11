package oAuth

import (
	"bytes"
	"encoding/json"
	"fmt"
	funcs "forum/funcs"
	"io"
	"log"
	"net/http"
	"net/url"
)

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		log.Println("Error: GET request for Google callback")
		funcs.ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}

	// Get the code from the request query param "code"
	code := r.URL.Query().Get("code")

	// Now we need to exchange this code for an access token
	accessToken, err := getGoogleAccessToken(code)
	if err != nil {
		funcs.ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	// Now we get the user information in exchange for the access token
	userInfo, err := getGoogleUserInfo(accessToken)
	if err != nil {
		funcs.ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	fmt.Println("User info:", userInfo)
}

func getGoogleAccessToken(code string) (string, error) {
	// This function is to get the access token in exchange for the code

	// Creating our request structure (it needs the clientID, clientSecret, code, and other necessary parameters)
	requestData := url.Values{}
	requestData.Set("client_id", googleClientID)
	requestData.Set("client_secret", googleClientSecret)
	requestData.Set("code", code)
	requestData.Set("redirect_uri", "http://localhost:8080/googlecallback") // Your callback URL
	requestData.Set("grant_type", "authorization_code")

	// Creating a new request
	req, err := http.NewRequest(
		"POST",
		"https://oauth2.googleapis.com/token",
		bytes.NewBufferString(requestData.Encode()),
	)
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request and get the response
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error getting response:", err)
		return "", err
	}
	defer res.Body.Close()

	// Check for HTTP status code
	if res.StatusCode != http.StatusOK {
		log.Println("Unexpected response status:", res.Status)
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response into a struct
	type GoogleTokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	var response GoogleTokenResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		log.Println("Error decoding response:", err)
		return "", err
	}

	if response.AccessToken == "" {
		log.Println("No access token in response")
		return "", fmt.Errorf("empty access token in response")
	}
	return response.AccessToken, nil
}

func getGoogleUserInfo(accessToken string) (string, error) {
	// This function is to get the user information in exchange for the access token

	// Creating a new request
	req, err := http.NewRequest(
		"GET",
		"https://www.googleapis.com/oauth2/v2/userinfo",
		nil,
	)
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}

	// Set the Authorization header
	authHeader := fmt.Sprintf("Bearer %s", accessToken)
	req.Header.Set("Authorization", authHeader)

	// Send the request and get the response
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error getting response:", err)
		return "", err
	}
	defer res.Body.Close()

	// Check for HTTP status code
	if res.StatusCode != http.StatusOK {
		log.Println("Unexpected response status:", res.Status)
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Read the response body and return it
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return "", err
	}

	return string(body), nil
}
