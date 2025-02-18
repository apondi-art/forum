package handlers

import (
	"html/template"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, errorMessage string, statusCode int) {
	w.WriteHeader(statusCode)
	temp, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, "Failed to load error page", http.StatusInternalServerError)
		return
	}
	temp.Execute(w, map[string]interface{}{
		"StatusCode":   statusCode,
		"ErrorMessage": errorMessage,
	})
}
