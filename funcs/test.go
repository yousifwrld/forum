package forum

import (
	"fmt"
)

func insertTestPost() error {
	// Insert users
	users := []struct {
		email, username, password string
	}{
		{"test@example.com", "FirstUser", "password"},
		{"test2@example.com", "SecondUser", "password"},
		{"test3@example.com", "ThirdUser", "password"},
		{"test4@example.com", "FourthUser", "password"},
		{"test5@example.com", "FifthUser", "password"},
	}

	for _, user := range users {
		_, err := database.Exec(`INSERT INTO user (email, username, password) VALUES (?, ?, ?)`, user.email, user.username, user.password)
		if err != nil {
			fmt.Printf("Error inserting user %s: %v\n", user.username, err)
		}
	}

	// Insert categories
	categories := []string{"⚽sports", "🎵music", "📰news", "🎮gaming", "🖥️technology"}

	for _, category := range categories {
		_, err := database.Exec(`INSERT INTO category (name) VALUES (?)`, category)
		if err != nil {
			fmt.Printf("Error adding category %s: %v\n", category, err)
		}
	}

	// Insert posts
	posts := []struct {
		title, content string
		userID         int
	}{
		{"First Test Post", "This is the content of the first test post.", 1},
		{"Second Test Post", "This is the content of the second test post.", 2},
		{"Third Test Post", "This is the content of the third test post.", 3},
		{"Fourth Test Post", "This is the content of the fourth test post.", 4},
		{"Fifth Test Post", "This is the content of the fifth test post.", 5},
	}

	for _, post := range posts {
		_, err := database.Exec(`INSERT INTO post (title, content, userID) VALUES (?, ?, ?)`, post.title, post.content, post.userID)
		if err != nil {
			fmt.Printf("Error inserting post: %v\n", err)
		}
	}

	// Insert post categories
	postCategories := []struct {
		postID, categoryID int
	}{
		{1, 1},
		{1, 2},
		{2, 1},
		{2, 3},
		{3, 4},
		{3, 5},
		{4, 1},
		{4, 2},
		{5, 3},
		{5, 4},
	}

	for _, pc := range postCategories {
		_, err := database.Exec(`INSERT INTO post_categories (postID, categoryID) VALUES (?, ?)`, pc.postID, pc.categoryID)
		if err != nil {
			fmt.Printf("Error inserting post category: %v\n", err)
		}
	}

	// Insert comments
	comments := []struct {
		userID, postID int
		comment        string
	}{
		{1, 1, "This is the first comment."},
		{2, 2, "This is the second comment."},
		{3, 3, "This is the third comment."},
		{4, 4, "This is the fourth comment."},
		{5, 5, "This is the fifth comment."},
		{1, 2, "Another comment on the second post."},
		{2, 3, "Another comment on the third post."},
		{3, 4, "Another comment on the fourth post."},
		{4, 5, "Another comment on the fifth post."},
		{5, 1, "Another comment on the first post."},
	}

	for _, comment := range comments {
		_, err := database.Exec(`INSERT INTO comment (userID, postID, comment) VALUES (?, ?, ?)`, comment.userID, comment.postID, comment.comment)
		if err != nil {
			fmt.Printf("Error inserting comment: %v\n", err)
		}
	}

	return nil
}
