package forum

import (
	"database/sql"
	"log"
)

func LikeComment(userID, commentID int) error {
	var liked, disliked bool

	tx, err := database.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Query the database to check if the comment is liked or disliked by the user
	err = tx.QueryRow(`SELECT is_like, is_dislike FROM comment_likes WHERE userID = ? AND commentID = ?`, userID, commentID).Scan(&liked, &disliked)
	if err != nil {
		if err == sql.ErrNoRows {
			// Insert a new like
			_, err = tx.Exec(`INSERT INTO comment_likes (userID, commentID, is_like, is_dislike) VALUES (?, ?, ?, ?)`, userID, commentID, true, false)
			if err != nil {
				log.Println("Error inserting comment like:", err)
				return err
			}
			log.Println("Inserted new comment like.")
		} else {
			log.Println("Error querying comment like:", err)
			return err
		}
	} else {
		// Update the like and dislike status
		_, err = tx.Exec(`UPDATE comment_likes SET is_like = ?, is_dislike = ? WHERE userID = ? AND commentID = ?`, !liked, false, userID, commentID)
		if err != nil {
			log.Println("Error updating comment like:", err)
			return err
		}
		log.Println("Updated comment like status.")
	}

	return nil
}

func DislikeComment(userID, commentID int) error {
	var liked, disliked bool

	tx, err := database.Begin()
	if (err != nil) {
		log.Println("Error starting transaction:", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Query the database to check if the comment is liked or disliked by the user
	err = tx.QueryRow(`SELECT is_like, is_dislike FROM comment_likes WHERE userID = ? AND commentID = ?`, userID, commentID).Scan(&liked, &disliked)
	if err != nil {
		if err == sql.ErrNoRows {
			// Insert a new dislike
			_, err = tx.Exec(`INSERT INTO comment_likes (userID, commentID, is_like, is_dislike) VALUES (?, ?, ?, ?)`, userID, commentID, false, true)
			if err != nil {
				log.Println("Error inserting comment dislike:", err)
				return err
			}
			log.Println("Inserted new comment dislike.")
		} else {
			log.Println("Error querying comment dislike:", err)
			return err
		}
	} else {
		// Update the like and dislike status
		_, err = tx.Exec(`UPDATE comment_likes SET is_like = ?, is_dislike = ? WHERE userID = ? AND commentID = ?`, false, !disliked, userID, commentID)
		if err != nil {
			log.Println("Error updating comment dislike:", err)
			return err
		}
		log.Println("Updated comment dislike status.")
	}

	return nil
}