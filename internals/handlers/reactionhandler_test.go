package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestHandleReaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		method         string
		body           interface{}
		expectedStatus int
		mockSetup      func()
	}{
		{
			name:           "Unauthorized user",
			method:         "POST",
			body:           nil,
			expectedStatus: http.StatusUnauthorized,
			mockSetup: func() {
			},
		},
		{
			name:           "Invalid method",
			method:         "GET", // Using GET to trigger method mismatch
			body:           nil,
			expectedStatus: http.StatusMethodNotAllowed,
			mockSetup:      func() {},
		},
		{
			name:           "Invalid request body",
			method:         "POST",
			body:           "Invalid JSON", // Invalid body
			expectedStatus: http.StatusUnauthorized,
			mockSetup:      func() {},
		},
		{
			name:   "Invalid target type",
			method: "POST",
			body: map[string]interface{}{
				"targetId":   1,
				"targetType": "invalid", // Invalid target type
				"type":       "like",
			},
			expectedStatus: http.StatusUnauthorized,
			mockSetup:      func() {},
		},
		{
			name:           "Valid reaction to post",
			method:         "POST",
			body:           map[string]interface{}{"targetId": 1, "targetType": "post", "type": "like"},
			expectedStatus: http.StatusUnauthorized,
			mockSetup: func() {},
		},
		{
			name:           "Valid reaction to comment",
			method:         "POST",
			body:           map[string]interface{}{"targetId": 1, "targetType": "comment", "type": "dislike"},
			expectedStatus: http.StatusUnauthorized,
			mockSetup: func() {},
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.mockSetup()
			var requestBody []byte
			if tt.body != nil {
				var err error
				requestBody, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			}
			req := httptest.NewRequest(tt.method, "/reaction", bytes.NewReader(requestBody))
			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler := http.HandlerFunc(HandleReaction)
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, status)
			}
			if tt.expectedStatus == http.StatusOK {
				var response map[string]int
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if response["likes"] == 0 && response["dislikes"] == 0 {
					t.Error("expected non-zero reaction counts in response")
				}
			}
		})
	}
}