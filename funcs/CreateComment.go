package forum

func CreateComment(userID, postID int, comment string) error {

	tx, err := database.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO comment (userID, postID, comment) VALUES (?, ?, ?)", userID, postID, comment)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
