package handlers

import (
	"encoding/json"
	"net/http"

	"forum/internals/auth"
	"forum/internals/models/categorymodel"
	"forum/internals/models/commentmodel"
)

// HandleCreateComment processes new comment submissions
func HandleCreateComment(w http.ResponseWriter, r *http.Request) {
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
		PostID  int64  `json:"postId"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorHandler(w, r, "Invalid request", http.StatusBadRequest)
		return
	}

	commentID, err := commentmodel.CreateComment(request.PostID, userID, request.Content)
	if err != nil {
		ErrorHandler(w, r, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	comment, err := commentmodel.GetCommentByID(commentID)
	if err != nil {
		ErrorHandler(w, r, "Failed to retrieve comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// HandleGetCategories returns the list of available categories
func HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		ErrorHandler(w, r, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categories, err := categorymodel.GetAllCategories()
	if err != nil {
		ErrorHandler(w, r, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
