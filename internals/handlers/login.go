package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"forum/internals/database"
	"forum/internals/models/usermodel"
)

// func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Serve the signup page (ensure "signup.html" exists inside "templates/")
// 	http.ServeFile(w, r, "templates/login.html")
// }

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Fatal(err)
		}
		if err = temp.Execute(w, nil); err != nil {
			fmt.Printf("error executing signup page : %v\n", err)
		}
	}

	if r.Method == http.MethodPost{
		// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		fmt.Println("Error parsing form:", err)
		return
	}


	// Retrieve form values
	email := r.FormValue("email")
	password := r.FormValue("password")

	authenticator, err := usermodel.AuthenticateUser(email, password)
	if err != nil || !authenticator {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		fmt.Println("error authenticating user:", err)
		return

	}

	newUser, err := usermodel.GetUserByEmail(email)
	if err != nil || newUser ==nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		fmt.Println("Error retrieving user:", err)
		return
	}

	// Create session
	session, err := usermodel.CreateSession(database.DB, newUser.ID)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.ID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // Set to true in production
		Path:     "/",
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		
	}

	

}
