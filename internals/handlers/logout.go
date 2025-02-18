package handlers

import (
	"net/http"
	"time"

	"forum/internals/auth"
	"forum/internals/database"
)

// LogoutHandler logs the user out by deleting their session
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session token from the cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// No session token found, just redirect to homepage
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Delete session from the database
	err = auth.DeleteSession(database.DB, cookie.Value)
	if err != nil {
		ErrorHandler(w, r, "Error logging out", http.StatusInternalServerError)
		return
	}

	// Expire the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Set expiration in the past
		HttpOnly: true,
		Path:     "/",
	})

	// Redirect to homepage after logout
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
