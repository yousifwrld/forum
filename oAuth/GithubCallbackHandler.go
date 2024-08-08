package oAuth

import (
	"bytes"
	"encoding/json"
	"fmt"
	funcs "forum/funcs"
	"io"
	"log"
	"net/http"
)

// this is the handler or path that github will send the request containing the code to
func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {

	//get the code from the request query param "code"
	code := r.URL.Query().Get("code")

	//now we need to exchange this code for an access token
	accessToken, err := getAccessToken(code)
	if err != nil {
		funcs.ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	//now we get the user information in exchange for the access token
	userInfo, err := getUserInfo(accessToken)
	if err != nil {
		funcs.ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	fmt.Println("User info:", userInfo)
}

func getAccessToken(code string) (string, error) {
	//this function is to get the access token in exchange of the code

	//creating our request structure (it needs the clientID, clientSecret and the code we received)
	request := struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
	}{
		//mapping the keys to the values
		ClientID:     githubClientID,
		ClientSecret: githubClientSecret,
		Code:         code,
	}

	reqJson, err := json.Marshal(request)
	if err != nil {
		log.Println("Error marshalling request:", err)
		return "", err
	}

	//creating a new request
	req, err := http.NewRequest( //could have used http.Post and made the process 10x easier but this gives us more control and we can set the needed headers
		"POST", // request method
		"https://github.com/login/oauth/access_token", //endpoint
		bytes.NewBuffer(reqJson),                      //request body
		//newbuffer creates the request body out of the byte slice
	)
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}
	//setting the request headers
	req.Header.Set("Content-Type", "application/json") // Content-Type set to JSON since our request body is in JSON format
	req.Header.Set("Accept", "application/json")       // Accept header set to JSON, indicating we prefer the response to be in JSON format

	//send the request and get the response
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error getting response:", err)
		return "", err
	}

	//close the response body
	defer res.Body.Close()

	// Check for HTTP status code
	if res.StatusCode != http.StatusOK {
		log.Println("Unexpected response status:", res.Status)
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	//since we set the Accept header to JSON, we need to decode the response into a struct
	type githubResponse struct {
		AccessToken string `json:"access_token"`
	}

	var response githubResponse
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

func getUserInfo(accessToken string) (string, error) {
	//this function is to get the user information in exchange of the access token

	//creating a new request
	req, err := http.NewRequest(
		"POST", // request method
		"https://api.github.com/user",
		nil, //no response body, because we need to send the token in the requeset header
	)

	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}

	// header should be in this format " Authorization: token {token} "

	//format and set the header
	AuthHeader := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", AuthHeader)

	//send the request and get the response
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error getting response:", err)
		return "", err
	}

	//close the response body
	defer res.Body.Close()

	// Check for HTTP status code
	if res.StatusCode != http.StatusOK {
		log.Println("Unexpected response status:", res.Status)
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	//read the response body and return it
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return "", err
	}

	return string(body), nil
}
