package forum

import (
	"fmt"
	"sort"
)

func GetPostComments(postID int) ([]Comment, error) {
	rows, err := database.Query(`
		SELECT u.username, c.commentID, c.comment, c.created_at, c.likes, c.dislikes
		FROM comment c
		JOIN user u ON c.userID = u.userID
		WHERE PostId = ?	
	`, postID)
	if err != nil {
		fmt.Println("hi")
		return nil, err
	}

	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment

		err := rows.Scan(&comment.Username, &comment.CommentID, &comment.Content, &comment.CreatedAt, &comment.Likes, &comment.Dislikes)
		if err != nil {
			return nil, err
		}
		comment.FormattedCreatedAt = comment.CreatedAt.Format("2006-01-02 15:04")
		comments = append(comments, comment)
	}

	// Check for errors during rows iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Sort comments by created_at
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.After(comments[j].CreatedAt)
	})

	return comments, nil
}
