package handlers

import (
	"html/template"
	"net/http"
)

type UserData struct {
	IsLoggedIn bool   // Indicates if the user is logged in
	Username   string // Username of the logged-in user
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Example: Check if the user is logged in (replace with actual authentication logic)
	isLoggedIn := true
	username := "JohnDoe" // Retrieve the username from the session or database

	// Prepare the data to pass to the template
	data := UserData{
		IsLoggedIn: isLoggedIn,
		Username:   username,
	}

	// Parse the template file
	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	// Execute the template with the data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
		return
	}
}
