package usermodel

import (
	"database/sql"
	"fmt"
	"log"

	"forum/internals/database"
)

type ReactionCounts struct {
	Likes    int
	Dislikes int
}

func HandlePostReaction(userID, postID int64, reactionType string) (ReactionCounts, error) {
	return getReactionData(userID, postID, 0, "post", reactionType)
}

func HandleCommentReaction(userID, commentID int64, reactionType string) (ReactionCounts, error) {
	return getReactionData(userID, 0, commentID, "comment", reactionType)
}

func getReactionData(userID, postID, commentID int64, targetType, reactionType string) (ReactionCounts, error) {
	var counts ReactionCounts

	// Use sql.NullInt64 to handle NULL values
	var postIDNull, commentIDNull sql.NullInt64
	if targetType == "post" {
		postIDNull = sql.NullInt64{Int64: postID, Valid: true}
		commentIDNull = sql.NullInt64{Valid: false} // comment_id is NULL
	} else if targetType == "comment" {
		commentIDNull = sql.NullInt64{Int64: commentID, Valid: true}
		postIDNull = sql.NullInt64{Valid: false} // post_id is NULL
	}

	// Check if the user has already reacted
	var existingReaction string
	err := database.DB.QueryRow(`
        SELECT reaction_type FROM likes_dislikes
        WHERE user_id = ? AND (post_id = ? OR comment_id = ?)
    `, userID, postIDNull, commentIDNull).Scan(&existingReaction)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Failed to query existing reaction: %v", err)
		return counts, fmt.Errorf("failed to query existing reaction: %v", err)
	}

	// Handle the reaction based on existing state
	if err := updateReaction(userID, postIDNull, commentIDNull, existingReaction, reactionType); err != nil {
		return counts, err
	}

	// Get updated counts
	return getUpdatedCounts(postIDNull, commentIDNull)
}

