package forum

import (
	"fmt"
	"forum/db"
	"log"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		DeleteCookiesAndSession(w, r)
		if r.URL.Query().Get("email") != "" || r.URL.Query().Get("username") != "" || r.URL.Query().Get("password") != "" {
			ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
			return
		}
		RenderTemplate(w, "templates/register.html", nil)
	case "POST":
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")

		err := validateCredentials(email, username, password)
		if err != nil {
			ErrorPages(w, r, "400", http.StatusBadRequest)
			return
		}

		exists, err := userExists(email, username)
		if err != nil {
			log.Println(err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}

		if exists {
			ErrorPages(w, r, "exists", http.StatusConflict, "templates/register.html")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

		if err != nil {
			fmt.Println(err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}
		userID, err := CreateUser(email, username, string(hashedPassword))
		if err != nil {
			log.Println(err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}

		//deletes any existing session and creates a new one
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

func CreateUser(email, username, password string) (int, error) {
	result, err := db.Database.Exec(`INSERT INTO USER (
			Email,
			Username,
			Password) VALUES (?, ?, ?)`, email, username, password)
	if err != nil {
		return 0, fmt.Errorf("500")
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("500")
	}

	users = append(users, User{UserID: int(userID), Email: email, Username: username, Password: password})
	return int(userID), nil
}

func userExists(email, username string) (bool, error) {
	var exists bool

	err := db.Database.QueryRow(`SELECT EXISTS(SELECT 1 FROM USER WHERE lower(Username) = ? OR Email = ?)`, strings.ToLower(username), email).Scan(&exists)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return exists, nil
}

func validateCredentials(email, username, password string) error {
	emailReg := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	reg := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)

	if !emailReg.MatchString(email) {
		return fmt.Errorf("invalid email")
	}

	if !reg.MatchString(username) {
		return fmt.Errorf("invalid username")
	}

	if len(password) < 8 || strings.TrimSpace(password) == "" {
		return fmt.Errorf("password too short")
	}
	return nil
}
