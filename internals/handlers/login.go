package handlers

import (
	"fmt"
	"forum/internals/models/usermodel"
	"net/http"
)

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Serve the signup page (ensure "signup.html" exists inside "templates/")
	http.ServeFile(w, r, "templates/login.html")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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
	
	email := r.FormValue("email")
	password := r.FormValue("password")

	authenticator,err :=usermodel.AuthenticateUser(email,password)
	if err !=nil || !authenticator{
		http.Error(w, "Error parsing form data", http.StatusNotFound)
		fmt.Println("error authenticating user:", err)
		return

	}

	// Success response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User logged in succesfully")

}
