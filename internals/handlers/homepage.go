package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"forum/internals/auth"
	"forum/internals/models/categorymodel"
)

// Homepage handles the main page
func Homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, r, "Page not found", http.StatusNotFound)
		return
	}

	// Get user session info
	userID, isLoggedIn := auth.GetUserFromSession(r)
	var userName string
	if isLoggedIn {
		userName, _ = auth.GetUserNameByID(userID)
	}

	// Get all categories for the sidebar
	categories, err := categorymodel.GetAllCategories()
	if err != nil {
		ErrorHandler(w, r, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	// Get category ID from query parameter
	categoryIDStr := r.URL.Query().Get("category")
	var categoryID int64
	if categoryIDStr != "" {
		categoryID, err = strconv.ParseInt(categoryIDStr, 10, 64)
		if err != nil {
			ErrorHandler(w, r, "Invalid category ID", http.StatusBadRequest)
			return
		}
	}

	// Check if showing liked posts
	showLiked := r.URL.Query().Get("liked") == "true"
	myPosts := r.URL.Query().Get("myPosts") == "true"

	// Get posts based on filters
	posts, err := categorymodel.GetPostsBySingleCategory(categoryID, userID, showLiked, myPosts)
	if err != nil {
		ErrorHandler(w, r, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := PageData{
		Posts:          posts,
		Categories:     categories,
		IsLoggedIn:     isLoggedIn,
		UserID:         userID,
		UserName:       userName,
		ActiveCategory: categoryID,
		ShowingLiked:   showLiked,
		MyPosts:        myPosts,
	}

	// Parse and execute template
	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		ErrorHandler(w, r, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		ErrorHandler(w, r, "Failed to execute template", http.StatusInternalServerError)
		return
	}
}
