package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	// Log raw query parameters
	// log.Printf("Raw query parameters: %v", r.URL.RawQuery) // Add this line

	// Parse category filter from query parameters
	categoryFilter := r.URL.Query().Get("categories")
	// log.Printf("Category filter string: %s", categoryFilter) // Add this line

	var categoryIDs []int64
	if categoryFilter != "" {
		categoryStrs := strings.Split(categoryFilter, ",")
		for _, categoryStr := range categoryStrs {
			categoryID, err := strconv.ParseInt(categoryStr, 10, 64)
			if err == nil {
				categoryIDs = append(categoryIDs, categoryID)
			} else {
				log.Printf("Error parsing category ID: %v", err) // Add this line
			}
		}
	}

	// log.Printf("Parsed category IDs: %v", categoryIDs) // Add this line

	// Display different content based on login status
	var posts []viewmodel.PostView
	var err error
	if len(categoryIDs) > 0 {
		// Fetch posts filtered by categories
		postIDs, err := categorymodel.GetPostsByCategories(categoryIDs)
		if err != nil {
			http.Error(w, "Failed to load posts", http.StatusInternalServerError)
			return
		}

		// log.Printf("Filtered post IDs: %v", postIDs) // Add this line

		posts, err = postmodel.GetPostsByIDs(postIDs)
	} else {
		if isLoggedIn {
			// Fetch posts normally for logged-in users
			posts, err = postmodel.GetFilteredPosts(userID, sql.NullInt64{Valid: false}, false)
		} else {
			// Fetch only public posts by passing a special filter (e.g., userID = 0 or -1)
			posts, err = postmodel.GetFilteredPosts(0, sql.NullInt64{Valid: false}, false)
		}
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
