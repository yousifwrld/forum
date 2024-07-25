package forum

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := getPosts()
	if err != nil {
		fmt.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	RenderTemplate(w, "templates/home.html", posts)

}
func getPosts() ([]Post, error) {
	rows, err := database.Query(`
        SELECT p.postID, p.title, p.content, p.created_at, u.username 
        FROM post p
        JOIN user u ON p.userID = u.userID
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Username); err != nil {
			return nil, err
		}
		post.FormattedCreatedAt = post.CreatedAt.Format("2006-01-02 15:04")

		// Fetch categories for this post
		catRows, err := database.Query(`
            SELECT c.categoryID, c.name 
            FROM category c
            JOIN post_categories pc ON c.categoryID = pc.categoryID
            WHERE pc.postID = ?
        `, post.ID)
		if err != nil {
			return nil, err
		}
		defer catRows.Close()

		var categories []Category
		for catRows.Next() {
			var category Category
			if err := catRows.Scan(&category.CategoryID, &category.Name); err != nil {
				return nil, err
			}
			categories = append(categories, category)
		}
		post.Categories = categories

		posts = append(posts, post)
	}
	return posts, nil
}
