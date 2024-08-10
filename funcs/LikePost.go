package forum

import (
	"database/sql"
	"forum/db"
	"log"
)

func LikePost(userID, postID int) (int, int, error) {
	var liked, disliked bool

	tx, err := db.Database.Begin()
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

	// Query the database to check if the post is liked or disliked by the user
	err = tx.QueryRow(`SELECT is_like, is_dislike FROM reaction WHERE userID = ? AND postID = ?`, userID, postID).Scan(&liked, &disliked)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = tx.Exec(`INSERT INTO reaction (userID, postID, is_like, is_dislike) VALUES (?, ?, ?, ?)`, userID, postID, true, false)
			if err != nil {
				log.Println("Error inserting reaction:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE post SET likes = likes + 1 WHERE postID = ?`, postID)
			if err != nil {
				log.Println("Error updating post likes:", err)
				return 0, 0, err
			}
			log.Println("Inserted new reaction and updated likes.")
		} else {
			log.Println("Error querying reaction:", err)
			return 0, 0, err
		}
	} else {
		if liked {
			_, err = tx.Exec(`UPDATE reaction SET is_like = ? WHERE userID = ? AND postID = ?`, false, userID, postID)
			if err != nil {
				log.Println("Error updating reaction:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE post SET likes = likes - 1 WHERE postID = ?`, postID)
			if err != nil {
				log.Println("Error updating post likes:", err)
				return 0, 0, err
			}
			log.Println("Removed like and updated likes.")
		} else if disliked {
			_, err = tx.Exec(`UPDATE reaction SET is_like = ?, is_dislike = ? WHERE userID = ? AND postID = ?`, true, false, userID, postID)
			if err != nil {
				log.Println("Error updating reaction:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE post SET likes = likes + 1, dislikes = dislikes - 1 WHERE postID = ?`, postID)
			if err != nil {
				log.Println("Error updating post likes and dislikes:", err)
				return 0, 0, err
			}
			log.Println("Changed dislike to like and updated counts.")
		} else {
			_, err = tx.Exec(`UPDATE reaction SET is_like = ? WHERE userID = ? AND postID = ?`, true, userID, postID)
			if err != nil {
				log.Println("Error updating reaction:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE post SET likes = likes + 1 WHERE postID = ?`, postID)
			if err != nil {
				log.Println("Error updating post likes and dislikes:", err)
				return 0, 0, err
			}
		}
	}

	// Retrieve the updated likes count
	var likesCount int
	var dislikesCount int
	err = tx.QueryRow(`SELECT likes, dislikes FROM post WHERE postID = ?`, postID).Scan(&likesCount, &dislikesCount)
	if err != nil {
		log.Println("Error querying post likes count:", err)
		return 0, 0, err
	}

	return likesCount, dislikesCount, nil
}

func DislikePost(userID, postID int) (int, int, error) {
	var liked, disliked bool

	tx, err := db.Database.Begin()
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

	// Query the database to check if the post is liked or disliked by the user
	err = tx.QueryRow(`SELECT is_like, is_dislike FROM reaction WHERE userID = ? AND postID = ?`, userID, postID).Scan(&liked, &disliked)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = tx.Exec(`INSERT INTO reaction (userID, postID, is_like, is_dislike) VALUES (?, ?, ?, ?)`, userID, postID, false, true)
			if err != nil {
				log.Println("Error inserting reaction:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE post SET dislikes = dislikes + 1 WHERE postID = ?`, postID)
			if err != nil {
				log.Println("Error updating post dislikes:", err)
				return 0, 0, err
			}
			log.Println("Inserted new reaction and updated dislikes.")
		} else {
			log.Println("Error querying reaction:", err)
			return 0, 0, err
		}
	} else {
		if disliked {
			_, err = tx.Exec(`UPDATE reaction SET is_dislike = ? WHERE userID = ? AND postID = ?`, false, userID, postID)
			if err != nil {
				log.Println("Error updating reaction:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE post SET dislikes = dislikes - 1 WHERE postID = ?`, postID)
			if err != nil {
				log.Println("Error updating post dislikes:", err)
				return 0, 0, err
			}
			log.Println("Removed dislike and updated dislikes.")
		} else if liked {
			_, err = tx.Exec(`UPDATE reaction SET is_like = ?, is_dislike = ? WHERE userID = ? AND postID = ?`, false, true, userID, postID)
			if err != nil {
				log.Println("Error updating reaction:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE post SET likes = likes - 1, dislikes = dislikes + 1 WHERE postID = ?`, postID)
			if err != nil {
				log.Println("Error updating post likes and dislikes:", err)
				return 0, 0, err
			}
			log.Println("Changed like to dislike and updated counts.")
		} else {
			_, err = tx.Exec(`UPDATE reaction SET is_dislike = ? WHERE userID = ? AND postID = ?`, true, userID, postID)
			if err != nil {
				log.Println("Error updating reaction:", err)
				return 0, 0, err
			}
			_, err = tx.Exec(`UPDATE post SET dislikes = dislikes + 1 WHERE postID = ?`, postID)
			if err != nil {
				log.Println("Error updating post dislikes:", err)
				return 0, 0, err
			}
		}
	}

	// Retrieve the updated dislikes count
	var dislikesCount int
	var likesCount int
	err = tx.QueryRow(`SELECT dislikes, likes FROM post WHERE postID = ?`, postID).Scan(&dislikesCount, &likesCount)
	if err != nil {
		log.Println("Error querying post dislikes count:", err)
		return 0, 0, err
	}

	return dislikesCount, likesCount, nil
}
