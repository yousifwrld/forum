package forum

import (
	"database/sql"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Query().Get("email") != "" || r.URL.Query().Get("username") != "" || r.URL.Query().Get("password") != "" {
			ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
			return
		}
		RenderTemplate(w, "templates/login.html", nil)
	case "POST":

		email := r.FormValue("email")
		password := r.FormValue("password")
		var hashedPassword string
		var userID int

		// get the hashed user password that matches the email
		err := database.QueryRow(`SELECT Password, UserID FROM USER WHERE Email = ?`, email).Scan(&hashedPassword, &userID)
		if err != nil {
			// if no results (0 rows found) means that the user doesn't exist
			if err == sql.ErrNoRows {
				log.Println(err)
				ErrorPages(w, r, "invalid email", http.StatusBadRequest, "templates/login.html")
				return
			} else {
				// any other error mainly related to the execution of the query
				log.Println(err)
				ErrorPages(w, r, "500", http.StatusInternalServerError)
				return
			}
		}
		// compare the stored hashed password with the login password, will return error if they don't match
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			log.Println(err)
			ErrorPages(w, r, "invalid password", http.StatusBadRequest, "templates/login.html")
			return
		}

		sessionID := SetCookies(w, r)
		database.Exec(`INSERT INTO SESSION (SessionID, UserID) VALUES (?, ?)`, sessionID, userID)

		http.Redirect(w, r, "/home", http.StatusFound)

	default:
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}
}
