package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"forum/internals/auth"
	"forum/internals/database"
	"forum/internals/handlers"
	"forum/internals/models/categorymodel"
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

	// Seed default categories
	err := categorymodel.SeedCategories()
	if err != nil {
		log.Fatalf("Failed to seed categories: %v", err)
	}

	// Start a goroutine to clean up expired sessions periodically
	go func() {
		for {
			if err := auth.CleanupExpiredSessions(database.DB); err != nil {
				log.Printf("Error cleaning up sessions: %v", err)
			}
			time.Sleep(1 * time.Hour)
		}
	}()

	http.HandleFunc("/", handlers.Homepage)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/signup", handlers.SignUpHandler)
	http.HandleFunc("/api/reaction", handlers.HandleReaction)
	http.HandleFunc("/api/comment", handlers.HandleCreateComment)
	http.HandleFunc("/api/posts/create", handlers.HandleCreatePost)
	http.HandleFunc("/api/categories", handlers.HandleGetCategories)
	http.HandleFunc("/dashboard", handlers.DashboardHandler)
	log.Println("Server listen on : http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
