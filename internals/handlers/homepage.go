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
		http.NotFound(w, r)
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
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	// Get category ID from query parameter
	categoryIDStr := r.URL.Query().Get("category")
	var categoryID int64
	if categoryIDStr != "" {
		categoryID, err = strconv.ParseInt(categoryIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
	}

	// Get posts for the selected category (or all posts if categoryID is 0)
	posts, err := categorymodel.GetPostsBySingleCategory(categoryID)
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
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
	}

	// Parse and execute template
	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
}
