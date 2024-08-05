package forum

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PostDetail struct {
	ID            int
	Title         string
	Content       string
	Username      string
	CreatedAt     string
	Likes         int
	Dislikes      int
	CommentsCount int
	Comments      []Comment
	Categories    []Category
}

func PostDetailHandler(w http.ResponseWriter, r *http.Request) {
	// extracting ID from the GET request
	postIDStr := strings.TrimPrefix(string(r.URL.Path), "/post/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ErrorPages(w, r, "404", http.StatusBadRequest)
		return
	}
	var post PostDetail
	var postCreatedAt time.Time
	var UserID int

	err = database.QueryRow(`
        SELECT postID, title, content, userID, created_at, likes, dislikes
        FROM post
        WHERE postID = ?`, postID).Scan(&post.ID, &post.Title, &post.Content, &UserID, &postCreatedAt, &post.Likes, &post.Dislikes)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			ErrorPages(w, r, "400", http.StatusBadRequest)
			return
		}
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
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

	err = database.QueryRow(`SELECT COUNT(*) FROM comment WHERE postID = ?`, postID).Scan(&post.CommentsCount)
	if err != nil {
		if err == sql.ErrNoRows {
			post.CommentsCount = 0
		} else {
			log.Println(err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}
	}

	post.CreatedAt = postCreatedAt.Format("2006-01-02 15:04")
	comments, err := GetPostComments(postID)
	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	post.Comments = comments

	categories, err := GetCategoriesForPost(postID)
	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}
	post.Categories = categories

	tmpl, err := template.New("post-detail.html").Funcs(template.FuncMap{
		"joinAndTrim": joinAndTrim,
	}).ParseFiles("templates/post-detail.html")

	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, post); err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}
}
