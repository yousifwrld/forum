package forum

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}
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
		categories, err := GetCategories(post.ID)
		if err != nil {
			return nil, err
		}
		// Get likes and dislikes for this post
		likes, dislikes, err := getLikesDislikes(post.ID)
		if err != nil {
			return nil, err
		}
		//save likes, dislikes and categories for each post
		post.Likes = likes
		post.Dislikes = dislikes
		post.Categories = categories

		// Append post to posts
		posts = append(posts, post)
	}
	// sort from latest to oldest based on created_at
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})
	return posts, nil
}

func getLikesDislikes(postID int) (int, int, error) {
	var likes, dislikes int
	err := database.QueryRow(`
		SELECT likes, dislikes
		FROM post
		WHERE postID = ?
	`, postID).Scan(&likes, &dislikes)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil
		}
		return 0, 0, err
	}
	return likes, dislikes, nil
}
