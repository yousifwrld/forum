package forum

import (
	"database/sql"
	"sync"
	"time"
)

// variables for openning one db connection only
var (
	database *sql.DB
	once     sync.Once
)

type Post struct {
	ID                 int       `json:"id"`
	Title              string    `json:"title"`
	Content            string    `json:"content"`
	CreatedAt          time.Time `json:"created_at"`
	FormattedCreatedAt string    `json:"formatted_created_at"`
	Username           string    `json:"username"`
	Categories         []string  `json:"categories"`
}

type Comment struct {
	Username           string    `json:"username"`
	Content            string    `json:"content"`
	CreatedAt          time.Time `json:"created_at"`
	FormattedCreatedAt string    `json:"formatted_created_at"`
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
