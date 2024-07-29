package forum

import (
	"database/sql"
	"log"
)

func LikeComment(userID, commentID int) (int, int, error) {
	var liked, disliked bool

	tx, err := database.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return 0, 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.QueryRow(`SELECT is_like, is_dislike FROM comment_likes WHERE userID = ? AND commentID = ?`, userID, commentID).Scan(&liked, &disliked)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = tx.Exec(`INSERT INTO comment_likes (userID, commentID, is_like, is_dislike) VALUES (?, ?, ?, ?)`, userID, commentID, true, false)
			if err != nil {
				log.Println("Error inserting comment like:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE comment SET likes = likes + 1 WHERE commentID = ?`, commentID)
			if err != nil {
				log.Println("Error updating comment likes:", err)
				return 0, 0, err
			}
			log.Println("Inserted new comment like and updated likes.")
		} else {
			log.Println("Error querying comment like:", err)
			return 0, 0, err
		}
	} else {
		if liked {
			_, err = tx.Exec(`UPDATE comment_likes SET is_like = ? WHERE userID = ? AND commentID = ?`, false, userID, commentID)
			if err != nil {
				log.Println("Error updating comment like:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE comment SET likes = likes - 1 WHERE commentID = ?`, commentID)
			if err != nil {
				log.Println("Error updating comment likes:", err)
				return 0, 0, err
			}
			log.Println("Removed like and updated likes.")
		} else if disliked {
			_, err = tx.Exec(`UPDATE comment_likes SET is_like = ?, is_dislike = ? WHERE userID = ? AND commentID = ?`, true, false, userID, commentID)
			if err != nil {
				log.Println("Error updating comment like:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE comment SET likes = likes + 1, dislikes = dislikes - 1 WHERE commentID = ?`, commentID)
			if err != nil {
				log.Println("Error updating comment likes and dislikes:", err)
				return 0, 0, err
			}
			log.Println("Changed dislike to like and updated counts.")
		} else {
			_, err = tx.Exec(`UPDATE comment_likes SET is_like = ? WHERE userID = ? AND commentID = ?`, true, userID, commentID)
			if err != nil {
				log.Println("Error updating comment like:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE comment SET likes = likes + 1 WHERE commentID = ?`, commentID)
			if err != nil {
				log.Println("Error updating comment likes:", err)
				return 0, 0, err
			}
		}
	}

	var likes, dislikes int
	err = tx.QueryRow(`SELECT likes, dislikes FROM comment WHERE commentID = ?`, commentID).Scan(&likes, &dislikes)
	if err != nil {
		log.Println("Error querying comment likes and dislikes count:", err)
		return 0, 0, err
	}

	return likes, dislikes, nil
}

func DislikeComment(userID, commentID int) (int, int, error) {
	var liked, disliked bool

	tx, err := database.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return 0, 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.QueryRow(`SELECT is_like, is_dislike FROM comment_likes WHERE userID = ? AND commentID = ?`, userID, commentID).Scan(&liked, &disliked)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = tx.Exec(`INSERT INTO comment_likes (userID, commentID, is_like, is_dislike) VALUES (?, ?, ?, ?)`, userID, commentID, false, true)
			if err != nil {
				log.Println("Error inserting comment dislike:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE comment SET dislikes = dislikes + 1 WHERE commentID = ?`, commentID)
			if err != nil {
				log.Println("Error updating comment dislikes:", err)
				return 0, 0, err
			}
			log.Println("Inserted new comment dislike and updated dislikes.")
		} else {
			log.Println("Error querying comment dislike:", err)
			return 0, 0, err
		}
	} else {
		if disliked {
			_, err = tx.Exec(`UPDATE comment_likes SET is_dislike = ? WHERE userID = ? AND commentID = ?`, false, userID, commentID)
			if err != nil {
				log.Println("Error updating comment dislike:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE comment SET dislikes = dislikes - 1 WHERE commentID = ?`, commentID)
			if err != nil {
				log.Println("Error updating comment dislikes:", err)
				return 0, 0, err
			}
			log.Println("Removed dislike and updated dislikes.")
		} else if liked {
			_, err = tx.Exec(`UPDATE comment_likes SET is_like = ?, is_dislike = ? WHERE userID = ? AND commentID = ?`, false, true, userID, commentID)
			if err != nil {
				log.Println("Error updating comment dislike:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE comment SET likes = likes - 1, dislikes = dislikes + 1 WHERE commentID = ?`, commentID)
			if err != nil {
				log.Println("Error updating comment likes and dislikes:", err)
				return 0, 0, err
			}
			log.Println("Changed like to dislike and updated counts.")
		} else {
			_, err = tx.Exec(`UPDATE comment_likes SET is_dislike = ? WHERE userID = ? AND commentID = ?`, true, userID, commentID)
			if err != nil {
				log.Println("Error updating comment dislike:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE comment SET dislikes = dislikes + 1 WHERE commentID = ?`, commentID)
			if err != nil {
				log.Println("Error updating comment dislikes:", err)
				return 0, 0, err
			}
		}
	}

	var dislikes, likes int
	err = tx.QueryRow(`SELECT dislikes, likes FROM comment WHERE commentID = ?`, commentID).Scan(&dislikes, &likes)
	if err != nil {
		log.Println("Error querying comment dislikes and likes count:", err)
		return 0, 0, err
	}

	return dislikes, likes, nil
}
