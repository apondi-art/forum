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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, isLoggedIn := auth.GetUserFromSession(r)
	if !isLoggedIn {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request struct {
		TargetID   int64  `json:"targetId"`   // ID of post or comment
		TargetType string `json:"targetType"` // "post" or "comment"
		Type       string `json:"type"`       // "like" or "dislike"
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Return updated counts
	var likes, dislikes int
	var err error
	switch request.TargetType {
	case "post":
		likes, dislikes, err = usermodel.HandlePostReaction(userID, request.TargetID, request.Type)
	case "comment":
		likes, dislikes, err = usermodel.HandleCommentReaction(userID, request.TargetID, request.Type)
	default:
		http.Error(w, "Invalid target type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to get reaction counts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"likes":    likes,
		"dislikes": dislikes,
	})
}
