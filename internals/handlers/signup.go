package handlers

import (
	"fmt"
	"net/http"

	"forum/internals/models/usermodel"
)

func SignupPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Serve the signup page (ensure "signup.html" exists inside "templates/")
	http.ServeFile(w, r, "templates/signup.html")
}

func SignUpHandlerProcess(w http.ResponseWriter, r *http.Request) {
	// Ensure it's a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		fmt.Println("Error parsing form:", err)
		return
	}

	fmt.Println("Received signup request")

	// Retrieve form values
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPass := r.FormValue("confirm_pass")

	// Validate password match
	if password != confirmPass {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := usermodel.PasswordHashing(password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		fmt.Println("Error hashing password:", err)
		return
	}
	fmt.Println(username, email, hashedPassword)
	// Store user in the database
	if err := usermodel.CreateUser(username, email, hashedPassword); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		fmt.Println("Error inserting user:", err)
		return
	}

	// Success response
	w.WriteHeader(http.StatusCreated)
	// fmt.Fprintln(w, "User created successfully")
}
