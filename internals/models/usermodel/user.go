package usermodel

import (
	"time"

	"forum/internals/database"
)

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

// Create a new user in the database
func CreateUser(username, email, password string) (int64, error) {
	query := `
        INSERT INTO Users (username, email, password)
        VALUES (?, ?, ?)
    `
	result, err := database.DB.Exec(query, username, email, password)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetUserByID retrieves a user by ID
func GetUserByID(id int64) (*User, error) {
	user := &User{}
	query := `SELECT id, username, email, password, created_at FROM Users WHERE id = ?`
	err := database.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, email, password, created_at FROM Users WHERE email = ?`
	err := database.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
