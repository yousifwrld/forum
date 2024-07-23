package forum

import (
	"log"
	"net/http"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Render the form
		RenderTemplate(w, "templates/create-post.html", nil)
	} else if r.Method == http.MethodPost {
		// Handle form submission
		title := r.FormValue("title")
		content := r.FormValue("content")

		if title == "" || content == "" {
			ErrorPages(w, r, "400", http.StatusBadRequest)
			return
		}

		// We should handle user authentication and get the userID here (sessions)
		userID := 1

		err := CreatePost(userID, title, content)
		if err != nil {
			log.Println(err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
	}
}

// TODO: return the categories in the function categories []int
func CreatePost(userID int, title, content string) error {
	tx, err := database.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO post (userID, title, content) VALUES (?, ?, ?)", userID, title, content)
	if err != nil {
		tx.Rollback()
		return err
	}

	// postID, err := result.LastInsertId()
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// for _, categoryID := range categories {
	// 	_, err := tx.Exec("INSERT INTO post_categories (postID, CategoryID) VALUES (?, ?)", postID, categoryID)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		return err
	// 	}
	// }

	return tx.Commit()
}
