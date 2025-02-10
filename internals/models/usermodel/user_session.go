package usermodel

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"
)

type Session struct {
	ID        string
	UserID    int64
	CreatedAt time.Time
	ExpiresAt time.Time
}

// CreateSession generates a new session for a user
func CreateSession(db *sql.DB, userID int64) (*Session, error) {
	// Generate random session ID
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("error generating session ID: %v", err)
	}
	sessionID := base64.URLEncoding.EncodeToString(b)

	// Create session with 24-hour expiration
	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Store in database
	query := `
        INSERT INTO Sessions (id, user_id, created_at, expires_at)
        VALUES (?, ?, ?, ?)
    `
	_, err := db.Exec(query, session.ID, session.UserID, session.CreatedAt, session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("error storing session: %v", err)
	}

	return session, nil
}

// GetSession retrieves a valid session from the database
func GetSession(db *sql.DB, sessionID string) (*Session, error) {
	session := &Session{}
	query := `
        SELECT id, user_id, created_at, expires_at
        FROM Sessions
        WHERE id = ? AND expires_at > datetime('now')
    `
	err := db.QueryRow(query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.CreatedAt,
		&session.ExpiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return session, nil
}
