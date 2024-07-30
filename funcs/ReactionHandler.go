package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}

	type request struct {
		ID     int    `json:"id"`
		Action int    `json:"action"`
		Type   string `json:"type"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(userIDKey).(int)
	fmt.Println("User ID from context:", userID)

	var likes, dislikes int
	var err error
	if req.Type == "post" {
		if req.Action == 1 {
			likes, dislikes, err = LikePost(userID, req.ID)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to like the post", http.StatusInternalServerError)
				return
			}
		} else if req.Action == 0 {
			dislikes, likes, err = DislikePost(userID, req.ID)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to dislike the post", http.StatusInternalServerError)
				return
			}
		}
	} else if req.Type == "comment" {
		if req.Action == 1 {
			likes, dislikes, err = LikeComment(userID, req.ID)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to like the comment", http.StatusInternalServerError)
				return
			}
		} else if req.Action == 0 {
			dislikes, likes, err = DislikeComment(userID, req.ID)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to dislike the comment", http.StatusInternalServerError)
				return
			}
		}
	}

	response := map[string]int{"dislikes": dislikes, "likes": likes}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
