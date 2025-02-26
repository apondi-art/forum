package commentmodel

import (
	"database/sql"
	"time"

	"forum/internals/database"
	"forum/internals/models/usermodel"
	"forum/internals/models/viewmodel"
)

type Comment struct {
	ID        int64
	PostID    int64
	UserID    int64
	Content   string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

func CreateComment(postID, userID int64, content string) (int64, error) {
	query := `INSERT INTO Comments (post_id, user_id, content) VALUES (?, ?, ?)`
	result, err := database.DB.Exec(query, postID, userID, content)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpdateComment updates an existing comment
func UpdateComment(commentID int64, content string) error {
	query := `
        UPDATE Comments
        SET content = ?, updated_at = CURRENT_TIMESTAMP
        WHERE id = ?
    `
	result, err := database.DB.Exec(query, content, commentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetCommentByID retrieves a single comment with all its details
func GetCommentByID(commentID int64) (*viewmodel.CommentView, error) {
	query := `
        SELECT c.id, c.content, u.username, c.created_at
        FROM Comments c
        JOIN Users u ON c.user_id = u.id
        WHERE c.id = ?
    `
	comment := &viewmodel.CommentView{}
	err := database.DB.QueryRow(query, commentID).Scan(
		&comment.ID,
		&comment.Content,
		&comment.Author,
		&comment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get likes/dislikes for the comment
	likes, dislikes, err := usermodel.GetCommentReactionCounts(comment.ID)
	if err != nil {
		return nil, err
	}
	comment.LikeCount = likes
	comment.DislikeCount = dislikes

	return comment, nil
}
