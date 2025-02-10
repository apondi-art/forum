package auth

import (
	"database/sql"
	"net/http"
	"time"

	"forum/internals/database"
	"forum/internals/models/usermodel"
)

// GetUserFromSession extracts user information from the session cookie
func GetUserFromSession(r *http.Request) (userID int64, isLoggedIn bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, false
	}

	session, err := usermodel.GetSession(database.DB, cookie.Value)
	if err != nil || session == nil {
		return 0, false
	}

	return session.UserID, true
}

// SetSessionCookie sets the session cookie in the response
func SetSessionCookie(w http.ResponseWriter, session *usermodel.Session) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   true, // Enable in production with HTTPS
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

// ClearSession removes a session from the database and clears the cookie
func ClearSession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		// Delete from database
		query := `DELETE FROM Sessions WHERE id = ?`
		_, err = database.DB.Exec(query, cookie.Value)
		if err != nil {
			return err
		}
	}

	// Clear cookie regardless of whether we found it in the DB
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	return nil
}

// CleanupExpiredSessions removes expired sessions from the database
func CleanupExpiredSessions(db *sql.DB) error {
	query := `DELETE FROM Sessions WHERE expires_at <= datetime('now')`
	_, err := db.Exec(query)
	return err
}
