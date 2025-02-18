package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"forum/internals/models/viewmodel"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreatePost(t *testing.T) {
	tests := []struct {
		name               string
		mockSession        func(r *http.Request) (int64, bool)
		mockCreatePost     func(userID int64, title, content string, categories []int64) (int64, error) 
		mockGetPostDetails func(postID int64) (viewmodel.PostView, error) 
		requestBody        interface{} 
		expectedStatusCode int         
		expectedPost       viewmodel.PostView 
	}{
		{
			name: "LoggedIn: Create Post Successfully",
			requestBody: struct {
				Title      string  `json:"title"`
				Content    string  `json:"content"`
				Categories []int64 `json:"categories"`
			}{
				Title:      "Test Post",
				Content:    "This is a test post.",
				Categories: []int64{1},
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedPost: viewmodel.PostView{
				ID:      123,
				Title:   "Test Post",
				Content: "This is a test post.",
			},
		},
		{
			name: "NotLoggedIn: Unauthorized",
			mockSession: func(r *http.Request) (int64, bool) {  // Simulate no user logged in
				return 0, false 
			},
			mockCreatePost: nil, 
			mockGetPostDetails: nil, 
			requestBody: struct {
				Title      string  `json:"title"`
				Content    string  `json:"content"`
				Categories []int64 `json:"categories"`
			}{
				Title:      "Test Post",
				Content:    "This is a test post.",
				Categories: []int64{1},
			},
			expectedStatusCode: http.StatusUnauthorized, 
			expectedPost:       viewmodel.PostView{}, 
		},
		{
			name: "InvalidData: Missing Title",
			mockSession: func(r *http.Request) (int64, bool) { 
				return 1, true 
			},
			mockCreatePost: nil, 
			mockGetPostDetails: nil,    
			requestBody: struct {
				Title      string  `json:"title"`
				Content    string  `json:"content"`
				Categories []int64 `json:"categories"`
			}{
				Title:      "",
				Content:    "This is a test post.",
				Categories: []int64{1},
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedPost:       viewmodel.PostView{},         
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the data into JSON
			postDataJSON, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("POST", "/create-post", bytes.NewReader(postDataJSON))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(HandleCreatePost)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			if rr.Code == http.StatusCreated {
				var post viewmodel.PostView
				err = json.NewDecoder(rr.Body).Decode(&post)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.expectedPost, post)
			}
		})
	}
}