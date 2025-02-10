package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"forum/internals/auth"
	"forum/internals/models/usermodel"
)

// HandleReaction processes post/comment likes and dislikes
func HandleReaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, isLoggedIn := auth.GetUserFromSession(r)
	if !isLoggedIn {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request struct {
		PostID    *int64 `json:"postId"`
		CommentID *int64 `json:"commentId"`
		Type      string `json:"type"` // "like" or "dislike"
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var postID, commentID sql.NullInt64
	if request.PostID != nil {
		postID = sql.NullInt64{Int64: *request.PostID, Valid: true}
	}
	if request.CommentID != nil {
		commentID = sql.NullInt64{Int64: *request.CommentID, Valid: true}
	}

	if err := usermodel.AddReaction(userID, postID, commentID, request.Type); err != nil {
		http.Error(w, "Failed to add reaction", http.StatusInternalServerError)
		return
	}

	// Return updated counts
	var likes, dislikes int
	var err error
	if postID.Valid {
		likes, dislikes, err = usermodel.GetReactionCounts(postID.Int64)
	} else {
		likes, dislikes, err = usermodel.GetCommentReactionCounts(commentID.Int64)
	}

	if err != nil {
		http.Error(w, "Failed to get reaction counts", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{
		"likes":    likes,
		"dislikes": dislikes,
	})
}
