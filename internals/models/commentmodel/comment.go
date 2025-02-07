package commentmodel

import (
	"database/sql"
	"time"

	"forum/internals/database"
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

func GetPostComments(postID int64) ([]Comment, error) {
	query := `
        SELECT id, post_id, user_id, content, created_at, updated_at
        FROM Comments
        WHERE post_id = ?
        ORDER BY created_at DESC
    `
	rows, err := database.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.UserID,
			&comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

