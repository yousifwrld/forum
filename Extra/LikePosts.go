package forum

func LikePost(userID, postID int, isLike bool) error {
	_, err := database.Exec("INSERT INTO LIKES (user_id, post_id, is_like) VALUES (?, ?, ?)", userID, postID, isLike)
	return err
}
