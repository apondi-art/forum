package postmodel

import (
	"database/sql"
	"errors"
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

func GetPostByID(id int64) (*Post, error) {
	post := &Post{}
	query := `
		SELECT id, user_id, title, content, created_at, updated_at 
		FROM Posts 
		WHERE id = ?
	`
	err := database.DB.QueryRow(query, id).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content,
		&post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func UpdatePost(postID int64, title, content string) error {
	query := `
        UPDATE Posts
        SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP
        WHERE id = ?
    `
	result, err := database.DB.Exec(query, title, content, postID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("post not found")
	}
	return nil
}
