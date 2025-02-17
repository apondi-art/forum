package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"forum/internals/auth"
	"forum/internals/models/categorymodel"
	"forum/internals/models/postmodel"
	"forum/internals/models/viewmodel"
)

type PageData struct {
	Posts          []viewmodel.PostView
	Categories     []categorymodel.Category
	IsLoggedIn     bool
	UserID         int64
	UserName       string
	ActiveCategory int64
	ShowingLiked   bool
}

// Template functions map
var funcMap = template.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format("Jan 02, 2006 15:04")
	},
	"len": func(v interface{}) int {
		if v == nil {
			return 0
		}
		switch val := v.(type) {
		case []viewmodel.CommentView:
			return len(val)
		default:
			return 0
		}
	},
}

func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is logged in
	userID, isLoggedIn := auth.GetUserFromSession(r)
	if !isLoggedIn {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var request struct {
		Title      string  `json:"title"`
		Content    string  `json:"content"`
		Categories []int64 `json:"categories"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create the post
	postID, err := postmodel.CreatePost(userID, request.Title, request.Content, request.Categories)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create post: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the created post
	post, err := viewmodel.GetPostWithDetails(postID)
	if err != nil {
		http.Error(w, "Failed to retrieve created post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}
