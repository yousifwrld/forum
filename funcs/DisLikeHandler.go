package forum

import (
	"encoding/json"
	"log"
	"net/http"
)

func DisLikeHandler(w http.ResponseWriter, r *http.Request) {
	// Define a struct to hold the incoming request data
	type request struct {
		PostID int `json:"postID"`
	}

	//decode the request body which contains the postID
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//dislike the post and get the dislikes and likes count
	dislikes, likes, err := DislikePost(1, req.PostID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to like the post", http.StatusInternalServerError)
		return
	}

	//store the dislikes and likes count in a map
	response := map[string]int{"dislikes": dislikes, "likes": likes}

	//set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	//encode the response and write it to the response writer to send it in json format
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
