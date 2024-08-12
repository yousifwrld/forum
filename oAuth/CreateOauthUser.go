package oAuth

import (
	"database/sql"
	"fmt"
	"forum/db"
	"log"
	"math/rand"
	"time"
)

// createUserByOAuth creates a new user with the provided OAuth information,
// linking the account appropriately based on existing entries in the database.
func createUserByOAuth(email, username, oauthProvider string, oauthUserID string) (int, error) {

	// Start a database transaction for atomicity
	tx, err := db.Database.Begin()
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %v", err)
	}

	// Ensure the transaction commits if all operations succeed, or rolls back if any operation fails
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Printf("Transaction failed and was rolled back due to panic: %v", p)
		} else if err != nil {
			tx.Rollback()
			log.Printf("Transaction rolled back due to error: %v", err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Printf("Error committing transaction: %v", err)
			}
		}
	}()

	// Check if the email already exists in the user table
	var existingUserID int
	err = tx.QueryRow(`SELECT userID FROM user WHERE email = ?`, email).Scan(&existingUserID)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("error checking for existing email: %v", err)
	}

	// If a user with the provided email already exists
	if existingUserID != 0 {
		// Check if there's already an OAuth link for this provider and userID
		var existingOAuthUserID string
		err = tx.QueryRow(`SELECT oauth_userID FROM user_oauth WHERE userID = ? AND oauth_provider = ?`, existingUserID, oauthProvider).Scan(&existingOAuthUserID)
		if err != nil && err != sql.ErrNoRows {
			return 0, fmt.Errorf("error checking for existing OAuth link: %v", err)
		}

		// If no OAuth link exists, insert a new one
		if err == sql.ErrNoRows {
			_, err = tx.Exec(
				`INSERT INTO user_oauth (userID, oauth_provider, oauth_userID) VALUES (?, ?, ?)`,
				existingUserID, oauthProvider, oauthUserID,
			)
			if err != nil {
				return 0, fmt.Errorf("error linking new OAuth information: %v", err)
			}

			// Clear the err variable after successful operation
			err = nil
		}

		// Return the existing user ID
		return existingUserID, nil
	}

	// If the email does not exist, generate a unique username
	originalUsername := username
	for {
		var existingUsername string
		err = tx.QueryRow(`SELECT username FROM user WHERE username = ?`, username).Scan(&existingUsername)

		if err == sql.ErrNoRows {
			break
		} else if err != nil {
			return 0, fmt.Errorf("error checking for duplicate username: %v", err)
		}

		// If the username is not unique, append a random number to it and try again
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomNumber := random.Intn(10000)
		username = fmt.Sprintf("%s%d", originalUsername, randomNumber)
	}

	// Insert the new user into the user table
	result, err := tx.Exec(
		`INSERT INTO user (email, username) VALUES (?, ?)`,
		email, username,
	)
	if err != nil {
		return 0, fmt.Errorf("error inserting new user: %v", err)
	}

	// Get the userID of the newly inserted user
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last inserted user ID: %v", err)
	}

	// Insert the OAuth information into the user_oauth table
	_, err = tx.Exec(
		`INSERT INTO user_oauth (userID, oauth_provider, oauth_userID) VALUES (?, ?, ?)`,
		userID, oauthProvider, oauthUserID,
	)
	if err != nil {
		return 0, fmt.Errorf("error inserting OAuth information: %v", err)
	}

	// The transaction will be committed by the deferred function if no error occurred
	return int(userID), nil
}
