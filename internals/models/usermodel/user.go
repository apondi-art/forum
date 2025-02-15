package usermodel

import (
	"fmt"
	"time"

	"forum/internals/database"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

// Create a new user in the database
func CreateUser(username, email, password string) error {
	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	_, err := database.DB.Exec(query, username, email, password)
	if err != nil {
		fmt.Printf("error inserting data %v\n", err)
		return err
	}
	return nil
}

func AuthenticateUser(email, password string) (bool, error) {
	var hashedPassword string

	// Fetch password hash from database
	err := database.DB.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&hashedPassword)
	if err != nil {
		return false, err // User not found
	}

	// Compare stored hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err // Incorrect password
	}

	return true, nil
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

func PasswordHashing(pasword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pasword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error occured during password hashing: %v\n", err)
		return "", err
	}
	return string(bytes), nil
}

func CompareHashedPassword(password string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

/*Declare a struct that holds the user  login credentials*/

// type LoginRequest struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// type LoginResponse struct {
// 	Message string `json:"message"`
// }

// type SignupRequest struct {
// 	Username    string `json:"username"`
// 	Email       string `json:"email"`
// 	Password    string `json:"password"`
// 	ConfirmPass string `json:"confirm_password"`
// }

// type SignupResponse struct {
// 	Message string `json:"message"`
// }
