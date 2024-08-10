package forum

import "forum/db"

func CreateComment(userID, postID int, comment string) error {

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
