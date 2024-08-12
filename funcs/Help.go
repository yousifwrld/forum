package forum

import (
	"forum/db"
	"net/http"
)

func Help(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}
	RenderTemplate(w, "templates/help.html", nil)

}
func Cr(userID, postID int, comment string) error {

	tx, err := db.Database.Begin()
	if err != nil {
		return err
	}

	// comment = strings.ReplaceAll(comment, "\r\n", "<br>")
	_, err = tx.Exec("INSERT INTO comment (userID, postID, comment) VALUES (?, ?, ?)", userID, postID, comment)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
