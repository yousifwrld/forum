package forum

import (
	"database/sql"
	"encoding/base64"
	"forum/db"
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
	image         []byte
	Base64Image   string
	Username      string
	CreatedAt     string
	Likes         int
	Dislikes      int
	CommentsCount int
	Comments      []Comment
	Categories    []Category
	UserLiked     bool
	UserDisliked  bool
}

func PostDetailHandler(w http.ResponseWriter, r *http.Request) {

	//get the sessionID from the cookie if available
	var err error
	var sessionID string
	sessionID, err = GetSessionFromCookie(r)
	if err != nil && err != http.ErrNoCookie {
		log.Println("Error getting session cookie:", err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	//if a session was available get the userID associated with the session
	reqUserID, err := GetIDFromSession(sessionID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error getting userID from session:", err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}
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

	//get the post from the database
	err = db.Database.QueryRow(`
        SELECT postID, title, content, image, userID, created_at, likes, dislikes
        FROM post
        WHERE postID = ?`, postID).Scan(&post.ID, &post.Title, &post.Content, &post.image, &UserID, &postCreatedAt, &post.Likes, &post.Dislikes)
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
	//get the username of the user that posted
	err = db.Database.QueryRow(`
		SELECT username 
		FROM user
		WHERE userID =?`, UserID).Scan(&post.Username)
	if err != nil {
		http.Error(w, "UserID not found", http.StatusNotFound)
		return
	}

	//get the count of comments for the post
	err = db.Database.QueryRow(`SELECT COUNT(*) FROM comment WHERE postID = ?`, postID).Scan(&post.CommentsCount)
	if err != nil {
		if err == sql.ErrNoRows {
			post.CommentsCount = 0
		} else {
			log.Println(err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}
	}

	//convert the image into a base64 string to display it in the HTML
	post.Base64Image = base64.StdEncoding.EncodeToString(post.image)

	//format the time into a more readable format
	post.CreatedAt = postCreatedAt.Format("2006-01-02 15:04")
	//get the comments for the post
	comments, err := GetPostComments(postID)
	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}
	post.Comments = comments

	//get the categories associated with the post
	categories, err := GetCategoriesForPost(postID)
	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}
	post.Categories = categories

	//get the user reactions for the post
	post.UserLiked, post.UserDisliked, err = getUserReactions(postID, reqUserID)
	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

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
