package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
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

	// Display different content based on login status
	var posts []viewmodel.PostView
	var err error
	if isLoggedIn {
		// Fetch posts normally for logged-in users
		posts, err = postmodel.GetFilteredPosts(userID, sql.NullInt64{Valid: false}, false)
	} else {
		// Fetch only public posts by passing a special filter (e.g., userID = 0 or -1)
		posts, err = postmodel.GetFilteredPosts(0, sql.NullInt64{Valid: false}, false)
	}
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
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
