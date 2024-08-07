package forum

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {

	// Check if the user is authenticated, and get the userID from the request context
	userID := r.Context().Value(userIDKey).(int)

	if r.Method == http.MethodGet {

		if r.URL.Query().Get("title") != "" || r.URL.Query().Get("content") != "" || r.URL.Query().Get("category") != "" {
			log.Println("Error: get request for create post")
			ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
			return
		}
		// Render the form with the categories to choose from

		//get them from the database
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
		// Render the template
		RenderTemplate(w, "templates/create-post.html", categories)
	} else if r.Method == http.MethodPost {

		// Parse the form, usin multipart form since thats the enctype we used in the html
		err := r.ParseMultipartForm(100 * 1024 * 1024) // 100MB max memory size for all the form data
		if err != nil {
			log.Println("Error parsing multipart form:", err)
			ErrorPages(w, r, "400", http.StatusBadRequest)
			return
		}

		// Get the categories from the form using r.MultipartForm.Value instead of r.FormValue it only returns the first value
		categories := r.MultipartForm.Value["category"]
		//get the form values normally
		title := r.FormValue("title")
		content := r.FormValue("content")

		// Get the image from the form
		//form file will return the first file provided
		image, _, err := r.FormFile("image")

		if err != nil {
			//if no image was uploaded which is okay because its optional
			if err == http.ErrMissingFile {
				// No file was uploaded
				log.Println("No image uploaded")
				image = nil // set it to nil since nothing was uploaded
			} else {
				// An error occurred while handling the uploaded file
				log.Println(err)
				ErrorPages(w, r, "500", http.StatusInternalServerError)
				return
			}
		} else {
			// Close the file when we're done
			defer image.Close()
		}

		//convert the image into bytes so that we can insert it into the database using blob data type which is used for big binary data
		var imageBytes []byte
		if image != nil {
			//read the image
			imageBytes, err = io.ReadAll(image)
			if err != nil {
				log.Println(err)
				ErrorPages(w, r, "500", http.StatusInternalServerError)
				return
			}
		}

		// handle No categories chosen
		if len(categories) == 0 {
			log.Println("No categories selected")
			ErrorPages(w, r, "400", http.StatusBadRequest)
			return
		}

		// handle title or content are empty
		if title == "" || content == "" {
			log.Println("Title or content is empty")
			ErrorPages(w, r, "400", http.StatusBadRequest)
			return
		}

		// handle image size is more than 20MB
		if len(imageBytes) > 20*1024*1024 {
			log.Println("Image size exceeds maximum size")
			ErrorPages(w, r, "400", http.StatusBadRequest)
			return
		}

		// Insert the post into the database
		err = CreatePost(userID, title, content, imageBytes, categories)
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

func CreatePost(userID int, title, content string, image []byte, categories []string) error {

	//begin a transaction
	tx, err := database.Begin()
	if err != nil {
		return err
	}

	//insert the post into the database with the image, it will be null incase of no image
	result, err := tx.Exec("INSERT INTO post (userID, title, content, image) VALUES (?, ?, ?, ?)", userID, title, content, image)
	if err != nil {
		tx.Rollback()
		return err
	}

	//get the post id to add the post categories
	postID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	var catIds []int

	//convert the ids to integers to insert them with the postID in the post_categories table
	for _, id := range categories {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("invalid category ID: %v", err)
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

	//commit the transaction incase of no errors
	return tx.Commit()
}
