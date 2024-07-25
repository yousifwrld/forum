package forum

import "fmt"

func insertTestPost() error {
	_, err := database.Exec(`
    INSERT INTO user (email, username, password) 
    VALUES (?, ?, ?)`,
		"test@example.com", "FirstUser", "password")
	if err != nil {
		fmt.Printf("Error inserting test post: %v\n", err)
	}
	_, err = database.Exec(`
    INSERT INTO user (email, username, password) 
    VALUES (?, ?, ?)`,
		"test@example.com", "SecondUser", "password")
	if err != nil {
		fmt.Printf("Error inserting test post: %v\n", err)
	}
	_, err = database.Exec(`
    INSERT INTO user (email, username, password) 
    VALUES (?, ?, ?)`,
		"test@example.com", "ThirdUser", "password")
	if err != nil {
		fmt.Printf("Error inserting test post: %v\n", err)
	}
	_, err = database.Exec(`
        INSERT INTO post (title, content, userID) 
        VALUES (?, ?, ?)`,
		"Test Post this is the first post for the program just for testing purposes not more than that I just want a little bit long pragraph for testing", "This is a test post content", 2)
	if err != nil {
		fmt.Printf("Error inserting test post: %v\n", err)
	}
	_, err = database.Exec(`
        INSERT INTO post (title, content, userID) 
        VALUES (?, ?, ?)`,
		"Test Post", "This is a test post content", 3)
	if err != nil {
		fmt.Printf("Error inserting test post: %v\n", err)
	}
	_, err = database.Exec(`
	INSERT INTO post (title, content, userID) 
	VALUES (?, ?, ?)`,
		"Test Post", "This is a test post content", 1)
	if err != nil {
		fmt.Printf("Error inserting test post: %v\n", err)
	}
	_, err = database.Exec(`
	INSERT INTO comment (userID, postID, comment) 
	VALUES (?, ?, ?)`,
		1, 2, "This the first comment ever")
	if err != nil {
		fmt.Printf("Error inserting test post: %v\n", err)
	}
	_, err = database.Exec(`
	INSERT INTO comment (userID, postID, comment) 
	VALUES (?, ?, ?)`,
		2, 3, "This the second comment and it is longer than the first one")
	if err != nil {
		fmt.Printf("Error inserting test post: %v\n", err)
	}
	_, err = database.Exec(`
	INSERT INTO comment (userID, postID, comment) 
	VALUES (?, ?, ?)`,
		3, 1, "This the second comment and iThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first onet is longer than the first one")
	if err != nil {
		fmt.Printf("Error inserting test post: %v\n", err)
	}
	_, err = database.Exec(`
	INSERT INTO comment (userID, postID, comment) 
	VALUES (?, ?, ?)`,
		3, 1, "This the second comment and iThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first onet is longer than the first one")
	if err != nil {
		fmt.Printf("Error inserting test post: %v\n", err)
	}
	_, err = database.Exec(`
	INSERT INTO comment (userID, postID, comment) 
	VALUES (?, ?, ?)`,
		3, 1, "This the second comment and iThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first onet is longer than the first one")
	if err != nil {
		fmt.Printf("Error inserting test post: %v\n", err)
	}
	_, err = database.Exec(`
	INSERT INTO comment (userID, postID, comment) 
	VALUES (?, ?, ?)`,
		3, 1, "This the second comment and iThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first onet is longer than the first one")
	if err != nil {
		fmt.Printf("Error inserting test post: %v\n", err)
	}
	_, err = database.Exec(`
	INSERT INTO comment (userID, postID, comment) 
	VALUES (?, ?, ?)`,
		3, 1, "This the second comment and iThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first onet is longer than the first one")
	if err != nil {
		fmt.Printf("Error inserting test post: %v\n", err)
	}
	_, err = database.Exec(`
	INSERT INTO comment (userID, postID, comment) 
	VALUES (?, ?, ?)`,
		3, 1, "This the second comment and iThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first oneThis the second comment and it is longer than the first onet is longer than the first one")
	if err != nil {
		fmt.Printf("Error inserting test post: %v\n", err)
	}
	_, err = database.Exec(`INSERT INTO category (name) VALUES (?)`, "sports")
	if err != nil {
		fmt.Printf("error adding category: %v", err)
	}

	_, err = database.Exec(`INSERT INTO category (name) VALUES (?)`, "music")
	if err != nil {
		fmt.Printf("error adding category: %v", err)
	}

	_, err = database.Exec(`INSERT INTO category (name) VALUES (?)`, "news")
	if err != nil {
		fmt.Printf("error adding category: %v", err)
	}

	return err
}
