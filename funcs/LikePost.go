package forum

import (
	"database/sql"
	"log"
)

func LikePost(userID, postID int) (int, error) {
	var liked, disliked bool

	// Query the database to check if the post is liked or disliked by the user
	err := database.QueryRow(`SELECT is_like, is_dislike FROM reaction WHERE userID = ? AND postID = ?`, userID, postID).Scan(&liked, &disliked)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = database.Exec(`INSERT INTO reaction (userID, postID, is_like, is_dislike) VALUES (?, ?, ?, ?)`, userID, postID, true, false)
			if err != nil {
				log.Println("Error inserting reaction:", err)
				return 0, err
			}
			_, err = database.Exec(`UPDATE post SET likes = likes + 1 WHERE postID = ?`, postID)
			if err != nil {
				log.Println("Error updating post likes:", err)
				return 0, err
			}
			log.Println("Inserted new reaction and updated likes.")
		} else {
			log.Println("Error querying reaction:", err)
			return 0, err
		}
	} else {
		if liked {
			_, err = database.Exec(`UPDATE reaction SET is_like = ? WHERE userID = ? AND postID = ?`, false, userID, postID)
			if err != nil {
				log.Println("Error updating reaction:", err)
				return 0, err
			}
			_, err = database.Exec(`UPDATE post SET likes = likes - 1 WHERE postID = ?`, postID)
			if err != nil {
				log.Println("Error updating post likes:", err)
				return 0, err
			}
			log.Println("Removed like and updated likes.")
		} else {
			_, err = database.Exec(`UPDATE reaction SET is_like = ? WHERE userID = ? AND postID = ?`, true, userID, postID)
			if err != nil {
				log.Println("Error updating reaction:", err)
				return 0, err
			}
			_, err = database.Exec(`UPDATE post SET likes = likes + 1 WHERE postID = ?`, postID)
			if err != nil {
				log.Println("Error updating post likes and dislikes:", err)
				return 0, err
			}
		}

		if disliked {
			_, err = database.Exec(`UPDATE reaction SET is_like = ?, is_dislike = ? WHERE userID = ? AND postID = ?`, true, false, userID, postID)
			if err != nil {
				log.Println("Error updating reaction:", err)
				return 0, err
			}
			_, err = database.Exec(`UPDATE post SET likes = likes + 1, dislikes = dislikes - 1 WHERE postID = ?`, postID)
			if err != nil {
				log.Println("Error updating post likes and dislikes:", err)
				return 0, err
			}
			log.Println("Changed dislike to like and updated counts.")
		} else {
			log.Println("Unhandled reaction state: liked =", liked, "disliked =", disliked)
		}
	}

	// Retrieve the updated likes count
	var likesCount int
	err = database.QueryRow(`SELECT likes FROM post WHERE postID = ?`, postID).Scan(&likesCount)
	if err != nil {
		log.Println("Error querying post likes count:", err)
		return 0, err
	}

	return likesCount, nil
}
