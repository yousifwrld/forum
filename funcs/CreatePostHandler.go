package forum

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {

	// Check if the user is authenticated, and get the userID from the request context
	userID := r.Context().Value(userIDKey).(int)

	if r.Method == http.MethodGet {
		// Render the form
		var categories []Category
		rows, err := database.Query(`SELECT categoryID, name FROM category`)
		if err != nil {
			log.Println("Error querying categories:", err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var category Category
			err := rows.Scan(&category.CategoryID, &category.Name)
			if err != nil {
				log.Println("Error scanning category:", err)
				ErrorPages(w, r, "500", http.StatusInternalServerError)
				return
			}
			categories = append(categories, category)
		}
		if err := rows.Err(); err != nil {
			log.Println("Error after iterating rows:", err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}
		RenderTemplate(w, "templates/create-post.html", categories)
	} else if r.Method == http.MethodPost {
		// Handle form submission
		err := r.ParseForm()
		if err != nil {
			ErrorPages(w, r, "400", http.StatusBadRequest)
		}
		categories := r.Form["category"]
		fmt.Println(categories)
		title := r.FormValue("title")
		content := r.FormValue("content")

		if title == "" || content == "" {
			ErrorPages(w, r, "400", http.StatusBadRequest)
			return
		}

		err = CreatePost(userID, title, content, categories)
		if err != nil {
			log.Println(err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
	}
}

// TODO: return the categories in the function categories []int
func CreatePost(userID int, title, content string, categories []string) error {
	tx, err := database.Begin()
	if err != nil {
		return err
	}

	result, err := tx.Exec("INSERT INTO post (userID, title, content) VALUES (?, ?, ?)", userID, title, content)
	if err != nil {
		tx.Rollback()
		return err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	var catIds []int

	for _, id := range categories {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return fmt.Errorf("400")
		}
		catIds = append(catIds, idInt)
	}

	for _, categoryID := range catIds {
		_, err := tx.Exec("INSERT INTO post_categories (postID, CategoryID) VALUES (?, ?)", postID, categoryID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
