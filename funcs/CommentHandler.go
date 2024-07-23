package forum

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("postID")
	if postIDStr == "" {
		http.Error(w, "Post ID is missing", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		// Render the form
		RenderTemplate(w, "templates/create-comment.html", postID)
	} else if r.Method == http.MethodPost {
		// Handle form submission
		content := r.FormValue("content")
		if content == "" {
			http.Error(w, "Content is missing", http.StatusBadRequest)
			return
		}
		// We should handle user authentication and get the userID here (sessions)
		sessionID, err := GetCookies(w, r)
		if err != nil {
			fmt.Println(err)
		}

		var userID int
		err = database.QueryRow(`
		SELECT userID
		FROM session
		WHERE sessionID=?`, sessionID).Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "fuck you", http.StatusBadRequest)
				fmt.Println(err)
				return
			} else {
				http.Error(w, "dfdf", http.StatusBadRequest)
				fmt.Println(err)
				return
			}
		}

		err = CreateComment(userID, postID, content)
		if err != nil {
			http.Error(w, "Unable to create comment", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home/post-detail/?postID="+postIDStr, http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
