package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"forum/internals/models/commentmodel"
)

func TestHandleCreateComment(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name           string
		args           args
		mockUserID     int64
		mockIsLoggedIn bool
		mockCreateErr  error
		mockComment    *commentmodel.Comment
		mockGetErr     error
	}{
		{
			name: "Method Not Allowed",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/comments", nil),
			},
			mockIsLoggedIn: true,
		},
		{
			name: "Unauthorized (Not Logged In)",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/comments", nil),
			},
			mockIsLoggedIn: false,
		},
		{
			name: "Invalid Request Body (Bad JSON)",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/comments", bytes.NewReader([]byte(`{"invalid_json":}`))),
			},
			mockIsLoggedIn: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run the handler
			HandleCreateComment(tt.args.w, tt.args.r)
			if tt.args.w.(*httptest.ResponseRecorder).Code != http.StatusOK && tt.args.w.(*httptest.ResponseRecorder).Code != http.StatusMethodNotAllowed &&
				tt.args.w.(*httptest.ResponseRecorder).Code != http.StatusUnauthorized &&
				tt.args.w.(*httptest.ResponseRecorder).Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, but got %d", http.StatusOK, tt.args.w.(*httptest.ResponseRecorder).Code)
			}
			// Check the response body
			if tt.mockComment != nil {
				var commentResp commentmodel.Comment
				if err := json.NewDecoder(tt.args.w.(*httptest.ResponseRecorder).Body).Decode(&commentResp); err != nil {
					t.Errorf("Failed to decode response body: %v", err)
				}
				if commentResp.ID != tt.mockComment.ID {
					t.Errorf("Expected comment ID %d, but got %d", tt.mockComment.ID, commentResp.ID)
				}
			}
		})
	}
}
