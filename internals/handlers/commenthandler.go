package handlers

import (
	"encoding/json"
	"net/http"

	"forum/internals/auth"
	"forum/internals/models/commentmodel"
)

// HandleCreateComment processes new comment submissions
func HandleCreateComment(w http.ResponseWriter, r *http.Request) {
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
		PostID  int64  `json:"postId"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	commentID, err := commentmodel.CreateComment(request.PostID, userID, request.Content)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	comment, err := commentmodel.GetCommentByID(commentID)
	if err != nil {
		http.Error(w, "Failed to retrieve comment", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comment)
}
