package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"forum/internals/models/usermodel"
)

// func SignupPageHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Serve the signup page (ensure "signup.html" exists inside "templates/")
// 	http.ServeFile(w, r, "templates/signup.html")
// }

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure it's a POST request
	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles("templates/signup.html")
		if err != nil {
			log.Fatal(err)
		}
		if err = temp.Execute(w, nil); err != nil {
			fmt.Printf("error executing signup page : %v\n", err)
		}
	}
	if r.Method == http.MethodPost {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			fmt.Println("Error parsing form:", err)
			return
		}

		

		// Retrieve form values
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPass := r.FormValue("confirm_pass")
		fmt.Println("Received signup request",username,email,password)


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

	// Redirect user to the dashboard after successful signup
	http.Redirect(w, r, "/login", http.StatusSeeOther)

	}

}
