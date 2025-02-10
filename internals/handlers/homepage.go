package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"

	"forum/internals/auth"
	"forum/internals/models/categorymodel"
	"forum/internals/models/postmodel"
	"forum/internals/models/viewmodel"
)

// Homepage handles the main page
func Homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Get user session info
	userID, isLoggedIn := auth.GetUserFromSession(r)

	// Get categories for filter
	categories, err1 := categorymodel.GetAllCategories()
	if err1 != nil {
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	// Parse filter parameters
	var selectedCategories []int64
	var showLiked bool

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		for _, catID := range r.Form["category"] {
			id, err := strconv.ParseInt(catID, 10, 64)
			if err == nil {
				selectedCategories = append(selectedCategories, id)
			}
		}

		showLiked = r.Form.Get("showLiked") == "on"
	}

	// Get filtered posts
	var posts []viewmodel.PostView
	var err error

	if len(selectedCategories) > 0 {
		for _, catID := range selectedCategories {
			categoryPosts, err := postmodel.GetFilteredPosts(userID, sql.NullInt64{Int64: catID, Valid: true}, showLiked)
			if err != nil {
				http.Error(w, "Failed to load posts", http.StatusInternalServerError)
				return
			}
			posts = append(posts, categoryPosts...)
		}
	} else {
		posts, err = postmodel.GetFilteredPosts(userID, sql.NullInt64{Valid: false}, showLiked)
		if err != nil {
			http.Error(w, "Failed to load posts", http.StatusInternalServerError)
			return
		}
	}

	// Prepare template data
	data := PageData{
		Posts:      posts,
		Categories: categories,
		IsLoggedIn: isLoggedIn,
		UserID:     userID,
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
