package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	a "forum/internals/models/usermodel"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Printf("Error during parsing template %v\n", err)
		return
	}
	tmpl.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the login page for GET requests
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Error loading login template", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error rendering login page", http.StatusInternalServerError)
		}
		return
	}

	// Handle POST request (login logic)
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Retrieve the form data
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Validate the login credentials
		if username == "" || password == "" {
			http.Error(w, "Both username and password are required", http.StatusBadRequest)
			return
		}

		// Here, you would typically check the credentials against a database or another storage system
		// For now, we'll just check if they match "admin" and "password"
		if username == "admin" && password == "password" {
			// Respond with a successful login message
			resp := a.LoginResponse{Message: "Login successful"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		} else {
			// Respond with an error message if login fails
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		}
		return
	}

	// If method is not GET or POST, return method not allowed
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the signup page for GET requests
	if r.Method == http.MethodGet {
		// Parse the signup template
		tmpl, err := template.ParseFiles("/home/hilaromondi/zero/forum/internals/templates/signup.html")
		if err != nil {
			http.Error(w, "Error loading signup template", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil) // Render the template with no data (just the empty form)
		if err != nil {
			http.Error(w, "Error rendering signup page", http.StatusInternalServerError)
			return
		}
		return
	}
	// Handle POST request (signup logic)
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Retrieve the form data
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPass := r.FormValue("confirm_pass")

		// Check if passwords match
		if password != confirmPass {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}

		// // Simple password validation (minimum length)
		// if len(password) < 8 {
		// 	http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		// 	return
		// }

		// Create a response struct
		if  username == "admin" && email == "d@gmail.com" &&  password == "password" {
			// Respond with a successful signup message
			resp := a.SignupResponse{Message: "Sign up successful"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		} else {
			// Respond with an error message if signup fails
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		}
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
