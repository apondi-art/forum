package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Printf("Error during parsing template %v\n", err)
		return
	}
	tmpl.Execute(w, nil)
}




