package usermodel

import (
	"database/sql"

	"forum/internals/database"
)

type Likes_Dislikes struct {
	ID           int64
	UserID       int64
	PostID       sql.NullInt64
	CommentID    sql.NullInt64
	ReactionType string
}

// Reaction related functions
func AddReaction(userID int64, postID, commentID sql.NullInt64, reactionType string) error {
	query := `
        INSERT INTO Likes_Dislikes (user_id, post_id, comment_id, reaction_type)
        VALUES (?, ?, ?, ?)
        ON CONFLICT (user_id, post_id, comment_id)
        DO UPDATE SET reaction_type = ?
    `
	_, err := database.DB.Exec(query, userID, postID, commentID, reactionType, reactionType)
	return err
}

func RemoveReaction(userID int64, postID, commentID sql.NullInt64) error {
	query := `
        DELETE FROM Likes_Dislikes
        WHERE user_id = ? AND post_id IS ? AND comment_id IS ?
    `
	_, err := database.DB.Exec(query, userID, postID, commentID)
	return err
}

func GetReactionCounts(postID int64) (likes int, dislikes int, error error) {
	query := `
        SELECT reaction_type, COUNT(*) 
        FROM Likes_Dislikes 
        WHERE post_id = ? AND comment_id IS NULL
        GROUP BY reaction_type
    `

	return populateReactions(query, postID)
}

// GetCommentReactionCounts gets likes/dislikes for a comment
func GetCommentReactionCounts(commentID int64) (likes int, dislikes int, err error) {
	query := `
        SELECT reaction_type, COUNT(*) 
        FROM Likes_Dislikes 
        WHERE comment_id = ? AND post_id IS NULL
        GROUP BY reaction_type
    `

	return populateReactions(query, commentID)
}

func populateReactions(query string, ID int64) (likes int, dislikes int, err error) {
	rows, err := database.DB.Query(query, ID)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var reactionType string
		var count int
		if err := rows.Scan(&reactionType, &count); err != nil {
			return 0, 0, err
		}
		if reactionType == "like" {
			likes = count
		} else {
			dislikes = count
		}
	}
	return likes, dislikes, nil
}
