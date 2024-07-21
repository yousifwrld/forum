package forum

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase() error {
	var err error
	// once ensures that the connection is only done once even if we call the function multiple times for some reason
	once.Do(func() {
		// Start database connection
		database, err = sql.Open("sqlite3", "forum.db")
		if err != nil {
			err = fmt.Errorf("error opening database: %v", err)
			return
		}
		// Create user table if it doesn't exist
		_, err = database.Exec(`CREATE TABLE IF NOT EXISTS USER (
			UserID INTEGER PRIMARY KEY AUTOINCREMENT,
			Email TEXT NOT NULL,
			Username TEXT NOT NULL,
			Password TEXT NOT NULL
		)`)
		if err != nil {
			err = fmt.Errorf("error creating user table: %v", err)
			return
		}

		_, err = database.Exec(`CREATE TABLE IF NOT EXISTS SESSION (
			SessionID TEXT PRIMARY KEY NOT NULL,
			UserID INTEGER NOT NULL,
			FOREIGN KEY (UserID) REFERENCES USER(UserID)
			)`)
		if err != nil {
			err = fmt.Errorf("error creating session table: %v", err)
			return
		}
	})
	return err
}
