package forum

func GetPostsByCategory(categoryID int) ([]Post, error) {
	rows, err := database.Query(`
        SELECT p.id, p.title, p.content, p.created_at, u.username
        FROM post p
        JOIN post_categories pc ON p.id = pc.post_id
        JOIN users u ON p.user_id = u.id
        WHERE pc.CategoryID = ?
    `, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Username)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
