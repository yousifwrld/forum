package forum

import (
	"encoding/json"
	"log"
	"net/http"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	// Define a struct to hold the incoming request data
	type request struct {
		PostID int `json:"postID"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	likes, err := LikePost(1, req.PostID) // Assume userID is 1 for now.
	if err != nil {
		log.Println(err)
		log.Println("dfdf")
		http.Error(w, "Failed to like the post", http.StatusInternalServerError)
		return
	}

	response := map[string]int{"likes": likes}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
