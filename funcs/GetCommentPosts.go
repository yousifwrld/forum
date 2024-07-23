package forum

import "time"

func GetCommentPosts(postID int) ([]Comment, error) {
	rows, err := database.Query(`
		SELECT u.username, c.comment, c.created_at
		FROM comment c
		JOIN user u ON c.userID = u.userID
		WHERE PostId = ?	
	`, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		var createdAt time.Time
		err := rows.Scan(&comment.Username, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comment.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		comments = append(comments, comment)
	}

	// Check for errors during rows iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
