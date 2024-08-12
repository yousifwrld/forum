package forum

import (
	"database/sql"
	"forum/db"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		DeleteCookiesAndSession(w, r)
		if r.URL.Query().Get("username") != "" || r.URL.Query().Get("password") != "" {
			ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
			return
		}
		RenderTemplate(w, "templates/login.html", nil)
	case "POST":

		username := strings.ToLower(strings.TrimSpace(r.PostFormValue("username")))
		password := strings.TrimSpace(r.FormValue("password"))

		if username == "" || password == "" {
			ErrorPages(w, r, "user not found", http.StatusBadRequest, "templates/login.html")
			return
		}
		var hashedPassword *string // Use a pointer to check for NULL
		var userID int

		// Get the hashed user password and user ID
		err := db.Database.QueryRow(`SELECT Password, UserID FROM USER WHERE Username = ?`, username).Scan(&hashedPassword, &userID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println(err)
				ErrorPages(w, r, "user not found", http.StatusBadRequest, "templates/login.html")
				return
			} else {
				log.Println(err)
				ErrorPages(w, r, "500", http.StatusInternalServerError)
				return
			}
		}

		// Check if the user has a password set
		if hashedPassword == nil {
			log.Println("Account created through OAuth, cannot log in through form.")
			ErrorPages(w, r, "oauth account", http.StatusForbidden, "templates/login.html")
			return
		}

		// Compare the stored hashed password with the login password
		err = bcrypt.CompareHashAndPassword([]byte(*hashedPassword), []byte(password))
		if err != nil {
			log.Println(err)
			ErrorPages(w, r, "invalid password", http.StatusBadRequest, "templates/login.html")
			return
		}

		// Deletes any existing session and creates a new one
		err = ManageSession(w, r, userID)
		if err != nil {
			log.Println(err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)

	default:
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}
}
