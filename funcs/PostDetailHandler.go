package forum

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

type PostDetail struct {
	ID        int
	Title     string
	Content   string
	Username  string
	CreatedAt string
	Comments  []Comment
}

func PostDetailHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("postID")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ErrorPages(w, r, "400", http.StatusBadRequest)
		return
	}
	var post PostDetail
	var postCreatedAt time.Time
	var UserID int

	err = database.QueryRow(`
        SELECT postID, title, content, userID, created_at
        FROM post
        WHERE postID = ?`, postID).Scan(&post.ID, &post.Title, &post.Content, &UserID, &postCreatedAt)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	err = database.QueryRow(`
		SELECT username 
		FROM user
		WHERE userID =?`, UserID).Scan(&post.Username)
	if err != nil {
		http.Error(w, "UserID not found", http.StatusNotFound)
		return
	}

	post.CreatedAt = postCreatedAt.Format("2006-01-02 15:04:05")
	comments, err := GetPostComments(postID)
	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	post.Comments = comments

	RenderTemplate(w, "templates/post-detail.html", post)
}
