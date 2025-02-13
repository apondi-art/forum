package auth

import (
	"database/sql"
	"errors"

	"forum/internals/database"
)

// GetUserNameByID retrieves the username by user ID
func GetUserNameByID(userID int64) (string, error) {
	var username string
	err := database.DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", err
	}
	return username, nil
}
