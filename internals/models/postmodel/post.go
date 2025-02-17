package postmodel

import (
	"database/sql"
	"time"

	"forum/internals/database"
)

type Post struct {
	ID        int64
	UserID    int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

func CreatePost(userID int64, title, content string, categoryIDs []int64) (int64, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Insert post
	query := `INSERT INTO Posts (user_id, title, content) VALUES (?, ?, ?)`
	result, err := tx.Exec(query, userID, title, content)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Add categories
	for _, categoryID := range categoryIDs {
		_, err = tx.Exec(`
            INSERT INTO Post_Categories (post_id, category_id)
            VALUES (?, ?)
        `, postID, categoryID)
		if err != nil {
			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return postID, nil
}
