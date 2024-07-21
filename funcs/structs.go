package forum

import (
	"database/sql"
	"sync"
)

// variables for openning one db connection only
var (
	database *sql.DB
	once     sync.Once
)

var users []User

type User struct {
	UserID                    int
	Email, Username, Password string
}
