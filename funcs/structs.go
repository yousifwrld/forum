package forum

import (
	"time"
)

type Post struct {
	ID                 int        `json:"id"`
	Title              string     `json:"title"`
	Content            string     `json:"content"`
	Image              []byte     `json:"image"`
	Base64Image        string     `json:"base64_image"`
	CreatedAt          time.Time  `json:"created_at"`
	FormattedCreatedAt string     `json:"formatted_created_at"`
	Username           string     `json:"username"`
	Categories         []Category `json:"categories"`
	Likes              int        `json:"likes"`
	Dislikes           int        `json:"dislikes"`
	FilterCategories   []Category
	Comments           int
	UserLiked          bool
	UserDisliked       bool
}

type Comment struct {
	CommentID          int       `json:"comment_id"`
	Username           string    `json:"username"`
	Content            string    `json:"content"`
	CreatedAt          time.Time `json:"created_at"`
	FormattedCreatedAt string    `json:"formatted_created_at"`
	Likes              int       `json:"likes"`
	Dislikes           int       `json:"dislikes"`
}

var users []User

type User struct {
	UserID                    int
	Email, Username, Password string
}

type Category struct {
	CategoryID int
	Name       string
}
