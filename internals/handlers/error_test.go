package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		errorMessage string
		statusCode   int
	}{
		{"Not Found", http.StatusNotFound},
		{"Internal Server Error", http.StatusInternalServerError},
	}

	// Mock template parsing
	templateParseFiles = func(filenames ...string) (*template.Template, error) {
		return template.New("error").Parse("{{.ErrorMessage}}")
	}

	for _, tt := range tests {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ErrorHandler(w, r, tt.errorMessage, tt.statusCode)
		})

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != tt.statusCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, tt.statusCode)
		}

		expected := tt.errorMessage
		body := rr.Body.String()
		if !strings.Contains(body, expected) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				body, expected)
		}
	}
}
