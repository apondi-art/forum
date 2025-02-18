package handlers

import (
	"encoding/json"
	"net/http"

	"forum/internals/auth"
	"forum/internals/models/usermodel"
)

// HandleReaction processes post/comment likes and dislikes
func HandleReaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		ErrorHandler(w, r, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, isLoggedIn := auth.GetUserFromSession(r)
	if !isLoggedIn {
		ErrorHandler(w, r, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request struct {
		TargetID   int64  `json:"targetId"`   // ID of post or comment
		TargetType string `json:"targetType"` // "post" or "comment"
		Type       string `json:"type"`       // "like" or "dislike"
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorHandler(w, r, "Invalid request", http.StatusBadRequest)
		return
	}

	// Return updated counts
	var counts usermodel.ReactionCounts
	var err error
	switch request.TargetType {
	case "post":
		counts, err = usermodel.HandlePostReaction(userID, request.TargetID, request.Type)
	case "comment":
		counts, err = usermodel.HandleCommentReaction(userID, request.TargetID, request.Type)
	default:
		ErrorHandler(w, r, "Invalid target type", http.StatusBadRequest)
		return
	}

	if err != nil {
		ErrorHandler(w, r, "Failed to get reaction counts", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"likes":    counts.Likes,
		"dislikes": counts.Dislikes,
	})
}
