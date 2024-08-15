package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	//check the request method
	if r.Method != http.MethodPost {
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}

	//struct for the requests
	type request struct {
		ID     int    `json:"id"`
		Action int    `json:"action"`
		Type   string `json:"type"`
	}

	//decode the request
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//get the user id from the request context
	userID := r.Context().Value(userIDKey).(int)
	fmt.Println("User ID from context:", userID)

	//call the appropriate function based on the request
	var likes, dislikes int
	var err error
	if req.Type == "post" { //if the request is for a post
		if req.Action == 1 { //if the user wants to like the post
			likes, dislikes, err = LikePost(userID, req.ID)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to like the post", http.StatusInternalServerError)
				return
			}
		} else if req.Action == 0 { //if the user wants to dislike the post
			dislikes, likes, err = DislikePost(userID, req.ID)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to dislike the post", http.StatusInternalServerError)
				return
			}
		}
	} else if req.Type == "comment" { //if the request is for a comment
		if req.Action == 1 { //if the user wants to like the comment
			likes, dislikes, err = LikeComment(userID, req.ID)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to like the comment", http.StatusInternalServerError)
				return
			}
		} else if req.Action == 0 { //if the user wants to dislike the comment
			dislikes, likes, err = DislikeComment(userID, req.ID)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to dislike the comment", http.StatusInternalServerError)
				return
			}
		}
	}

	//get the user's reaction for the post
	var userLiked, userDisliked bool
	if req.Type == "post" {
		userLiked, userDisliked, err = getUserReactions(req.ID, userID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to get user reactions", http.StatusInternalServerError)
			return
		}
	}

	fmt.Println("Likes:", likes, "Dislikes:", dislikes, "User liked:", userLiked, "User disliked:", userDisliked)

	//create the response
	response := map[string]interface{}{
		"dislikes":      dislikes,
		"likes":         likes,
		"user_disliked": userDisliked,
		"user_liked":    userLiked,
	}
	//encode the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
