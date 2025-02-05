package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"forum/internals/handlers"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Println("Too much arguments")
		return
	}

	http.HandleFunc("/", handlers.Homepage)
	fmt.Println("Server listen on : http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
