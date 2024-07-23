package forum

import (
	"log"
	"net/http"
	"strconv"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("postID")
	if postIDStr == "" {
		ErrorPages(w, r, "400", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ErrorPages(w, r, "400", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		// Render the form
		RenderTemplate(w, "templates/create-comment.html", postID)
	} else if r.Method == http.MethodPost {
		// Handle form submission
		content := r.FormValue("content")
		if content == "" {
			ErrorPages(w, r, "400", http.StatusBadRequest)
			return
		}
		// We should handle user authentication and get the userID here (sessions)
		// will change the logic and get the userID from the request context, and will authenticate using middleware function

		userID := 1

		err = CreateComment(userID, postID, content)
		if err != nil {
			log.Println(err)
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home/post-detail/?postID="+postIDStr, http.StatusFound)
	} else {
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
	}
}
