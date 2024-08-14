package oAuth

import (
	"database/sql"
	"fmt"
	"forum/db"
	"log"
	"math/rand"
	"strings"
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
	var emailExists bool

	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM user WHERE email = ?)`, email).Scan(&emailExists)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("error checking for existing email: %v", err)
	}

	if emailExists {
		return 0, fmt.Errorf("email already exists")
	}

	// if the email does not exist, create a new user with a unique username
	originalUsername := username

	//infinite loop until we find a unique username
	for {
		var existingUsername string
		err = tx.QueryRow(`SELECT username FROM user WHERE lower(username) = ?`, strings.ToLower(username)).Scan(&existingUsername)

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
		`INSERT INTO user (email, username, oauth_provider, oauth_userID) VALUES (?, ?, ?, ?)`,
		email, username, oauthProvider, oauthUserID,
	)
	if err != nil {
		return 0, fmt.Errorf("error inserting new user: %v", err)
	}

	// Get the userID of the newly inserted user
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last inserted user ID: %v", err)
	}

	// The transaction will be committed by the deferred function if no error occurred
	return int(userID), nil
}
