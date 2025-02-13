package handlers

import (
	"html/template"
	"net/http"
	"time"

	"forum/internals/database"
	"forum/internals/models/usermodel"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Failed to load login page", http.StatusInternalServerError)
			return
		}
		temp.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// Retrieve form values
		email := r.FormValue("email")
		password := r.FormValue("password")

		authenticated, err := usermodel.AuthenticateUser(email, password)
		if err != nil || !authenticated {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Get user info
		newUser, err := usermodel.GetUserByEmail(email)
		if err != nil || newUser == nil {
			http.Error(w, "User not found", http.StatusInternalServerError)
			return
		}

		// Create session
		session, err := usermodel.CreateSession(database.DB, newUser.ID)
		if err != nil {
			http.Error(w, "Error creating session", http.StatusInternalServerError)
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
