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
		posts = append(posts, post)
	}
	return posts, nil
}
