package forum

import (
	"net/http"
	"text/template"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := getPosts()
	if err != nil {
		http.Error(w, "Unable to retrieve posts", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, posts)
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
