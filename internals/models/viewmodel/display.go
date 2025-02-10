package viewmodel

import (
	"database/sql"
	"time"

	"forum/internals/database"
	"forum/internals/models/usermodel"
)

// PostView represents all data needed to display a post
type PostView struct {
	ID           int64
	Title        string
	Content      string
	Author       string
	CreatedAt    time.Time
	UpdatedAt    sql.NullTime
	Categories   []string
	LikeCount    int
	DislikeCount int
	Comments     []CommentView
}

// CommentView represents all data needed to display a comment
type CommentView struct {
	ID           int64
	Content      string
	Author       string
	CreatedAt    time.Time
	LikeCount    int
	DislikeCount int
}

// GetPostWithDetails retrieves a post with all its associated data
func GetPostWithDetails(postID int64) (*PostView, error) {
	// Get basic post info
	query := `
        SELECT p.id, p.title, p.content, u.username, p.created_at, p.updated_at
        FROM Posts p
        JOIN Users u ON p.user_id = u.id
        WHERE p.id = ?
    `
	post := &PostView{}
	err := database.DB.QueryRow(query, postID).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Author,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get categories
	categoryQuery := `
        SELECT c.name
        FROM Categories c
        JOIN Post_Categories pc ON c.id = pc.category_id
        WHERE pc.post_id = ?
    `
	rows, err := database.DB.Query(categoryQuery, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		post.Categories = append(post.Categories, category)
	}

	// Get likes/dislikes count
	post.LikeCount, post.DislikeCount, err = usermodel.GetReactionCounts(postID)
	if err != nil {
		return nil, err
	}

	// Get comments
	post.Comments, err = GetPostComments(postID)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// GetPostComments retrieves all comments for a post with author and reaction info
func GetPostComments(postID int64) ([]CommentView, error) {
	query := `
        SELECT c.id, c.content, u.username, c.created_at
        FROM Comments c
        JOIN Users u ON c.user_id = u.id
        WHERE c.post_id = ?
        ORDER BY c.created_at DESC
    `
	rows, err := database.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []CommentView
	for rows.Next() {
		var comment CommentView
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.Author,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get likes/dislikes for comment
		likes, dislikes, err := usermodel.GetCommentReactionCounts(comment.ID)
		if err != nil {
			return nil, err
		}
		comment.LikeCount = likes
		comment.DislikeCount = dislikes

		comments = append(comments, comment)
	}

	return comments, nil
}
