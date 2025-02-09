package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"forum/internals/database"
	"forum/internals/handlers"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Println("Too much arguments")
		return
	}
	
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	http.HandleFunc("/", handlers.Homepage)
   	http.HandleFunc("/signup", handlers.SignupPageHandler)
	http.HandleFunc("/login",handlers.LoginPageHandler)
	http.HandleFunc("/login/dashboard",handlers.LoginHandler)
	http.HandleFunc("/signup/process", handlers.SignUpHandlerProcess)
	log.Println("Server listen on : http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
