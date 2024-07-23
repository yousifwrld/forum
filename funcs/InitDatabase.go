package forum

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase() error {
	var err error
	once.Do(func() {
		// Start database connection
		database, err = sql.Open("sqlite3", "forum.db")
		if err != nil {
			err = fmt.Errorf("error opening database: %v", err)
			return
		}
		// Create user table if it doesn't exist
		_, err = database.Exec(`CREATE TABLE IF NOT EXISTS user (
			userID INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL,
			username TEXT NOT NULL,
			password TEXT NOT NULL
		)`)
		if err != nil {
			err = fmt.Errorf("error creating user table: %v", err)
			return
		}

		_, err = database.Exec(`CREATE TABLE IF NOT EXISTS session (
			sessionID TEXT PRIMARY KEY NOT NULL,
			userID INTEGER NOT NULL,
			FOREIGN KEY (userID) REFERENCES USER(userID)
		)`)
		if err != nil {
			err = fmt.Errorf("error creating session table: %v", err)
			return
		}

		// Create posts table if it doesn't exist
		_, err = database.Exec(`CREATE TABLE IF NOT EXISTS post (
			postID INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			userID INTEGER,
			FOREIGN KEY (userID) REFERENCES USER(userID)
		)`)
		if err != nil {
			err = fmt.Errorf("error creating post table: %v", err)
			return
		}

		// Create post_categories table if it doesn't exist
		_, err = database.Exec(`CREATE TABLE IF NOT EXISTS post_categories (
			postID INTEGER NOT NULL,
			categoryID INTEGER NOT NULL,
			PRIMARY KEY (postID, categoryID),
			FOREIGN KEY (postID) REFERENCES POST(postID),
			FOREIGN KEY (categoryID) REFERENCES CATEGORY(categoryID)
		)`)
		if err != nil {
			err = fmt.Errorf("error creating post_categories table: %v", err)
			return
		}

		// Create likes table if it doesn't exist
		_, err = database.Exec(`CREATE TABLE IF NOT EXISTS like (
			userID INTEGER NOT NULL,
			postID INTEGER NOT NULL,
			is_like BOOLEAN,
			is_dislike BOOLEAN,
			PRIMARY KEY (userID, postID),
			FOREIGN KEY (userID) REFERENCES USER(userID),
			FOREIGN KEY (postID) REFERENCES POST(postID)
		)`)
		if err != nil {
			err = fmt.Errorf("error creating likes table: %v", err)
			return
		}

		// Create comments table if it doesn't exist
		_, err = database.Exec(`CREATE TABLE IF NOT EXISTS comment (
			commentID INTEGER PRIMARY KEY AUTOINCREMENT,
			userID INTEGER NOT NULL,
			postID INTEGER NOT NULL,
			comment TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (userID) REFERENCES USER(userID),
			FOREIGN KEY (postID) REFERENCES POST(postID)
		)`)
		if err != nil {
			err = fmt.Errorf("error creating comments table: %v", err)
			return
		}
		// Insert a test post
		err = insertTestPost()
		if err != nil {
			fmt.Printf("Error inserting test post: %v\n", err)
		}
	})
	return err
}