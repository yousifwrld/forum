package forum

func GetCategoriesForPost(postID int) ([]Category, error) {
	// Fetch categories for this post
	catRows, err := database.Query(`
	SELECT c.categoryID, c.name 
	FROM category c
	JOIN post_categories pc ON c.categoryID = pc.categoryID
	WHERE pc.postID = ?
`, postID)

	if err != nil {
		return nil, err
	}
	defer catRows.Close()

	var categories []Category
	for catRows.Next() {
		var category Category
		if err := catRows.Scan(&category.CategoryID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
