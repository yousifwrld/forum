package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}

	type request struct {
		CommentID int `json:"commentID"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Println(req)

	userID := r.Context().Value(userIDKey).(int)
	fmt.Println("User ID from context:", userID)

	likes, dislikes, err := LikeComment(userID, req.CommentID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to like the comment", http.StatusInternalServerError)
		return
	}

	response := map[string]int{"dislikes": dislikes, "likes": likes}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func DisLikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}

	type request struct {
		CommentID int `json:"commentID"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	fmt.Println(req)

	userID := r.Context().Value(userIDKey).(int)
	fmt.Println("User ID from context:", userID)

	dislikes, likes, err := DislikeComment(userID, req.CommentID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to dislike the comment", http.StatusInternalServerError)
		return
	}

	response := map[string]int{"dislikes": dislikes, "likes": likes}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
