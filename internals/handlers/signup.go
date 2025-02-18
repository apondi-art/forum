package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"forum/internals/models/usermodel"
)


type SignupData struct {
	Error string
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	temp, err := template.ParseFiles("templates/signup.html")
    if err != nil {
        log.Fatal(err)
    }

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
			ErrorHandler(w, r, "Error parsing form data", http.StatusBadRequest)
			fmt.Println("Error parsing form:", err)
			return
		}

		// Retrieve form values
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPass := r.FormValue("confirm_pass")
         
         // Validate password match
		if password != confirmPass {
			data := SignupData{
				Error: "Passwords do not match",
			}
			temp.Execute(w, data)
			return
		}
		// Validate password match
		if password != confirmPass {
			ErrorHandler(w, r, "Passwords do not match", http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := usermodel.PasswordHashing(password)
		if err != nil {
			ErrorHandler(w, r, "Error hashing password", http.StatusInternalServerError)
			fmt.Println("Error hashing password:", err)
			return
		}

		// Store user in the database
		if err := usermodel.CreateUser(username, email, hashedPassword); err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			fmt.Println("Error inserting user:", err)
			return
		}

		// Redirect user to the login after successful signup
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	}
}
