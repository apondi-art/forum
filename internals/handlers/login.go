package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"forum/internals/auth"
	"forum/internals/database"
	"forum/internals/models/usermodel"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles("templates/login.html")
		if err != nil {
			ErrorHandler(w, r, "Failed to load login page", http.StatusInternalServerError)
			return
		}
		temp.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			ErrorHandler(w, r, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// Retrieve form values
		email := r.FormValue("email")
		password := r.FormValue("password")

		authenticated, err := usermodel.AuthenticateUser(email, password)
		if err != nil || !authenticated {
			temp, err := template.ParseFiles("templates/login.html")
			if err != nil {
				ErrorHandler(w, r, "Failed to load login page", http.StatusInternalServerError)
				return
			}
			temp.Execute(w, map[string]interface{}{
				"ErrorMessage": "Please check your email or password, then try again !!!",
			})
			return
		}

		// Get user info
		newUser, err := usermodel.GetUserByEmail(email)
		if err != nil || newUser == nil {
			ErrorHandler(w, r, "User not found", http.StatusInternalServerError)
			return
		}
		oldsession, err := usermodel.GetSessionbyUserID(database.DB, newUser.ID)
		fmt.Println(oldsession)
		if err != nil {
			ErrorHandler(w, r, "Error getting session", http.StatusInternalServerError)
			return
		} else {
			if oldsession != nil {
				auth.DeleteSession(database.DB, oldsession.ID)
			}
		}
		// Create session
		session, err := usermodel.CreateSession(database.DB, newUser.ID)
		if err != nil {
			ErrorHandler(w, r, "Error creating session", http.StatusInternalServerError)
			return
		}

		// Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token", // Ensure it matches GetUserFromSession
			Value:    session.ID,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   false, // Set to true in production with HTTPS
			Path:     "/",
		})

		// Redirect to homepage
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
