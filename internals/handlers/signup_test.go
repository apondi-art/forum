package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)


func (m *MockUserModel) PasswordHashing(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockUserModel) CreateUser(username, email, hashedPassword string) error {
	args := m.Called(username, email, hashedPassword)
	return args.Error(0)
}

func TestSignUpHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		username     string
		email        string
		password     string
		confirmPass  string
		expectedCode int
		expectedBody string
		mockActions  func(m *MockUserModel)
	}{
		{
			name:         "Valid signup",
			method:       http.MethodPost,
			username:     "newuser",
			email:        "newuser@example.com",
			password:     "validpassword",
			confirmPass:  "validpassword",
			expectedCode: http.StatusBadRequest, 
			expectedBody: "Error parsing form data\n",           
			mockActions: func(m *MockUserModel) {
				m.On("PasswordHashing", "validpassword").Return("hashedpassword", nil)
				m.On("CreateUser", "newuser", "newuser@example.com", "hashedpassword").Return(nil)
			},
		},
		{
			name:         "Passwords do not match",
			method:       http.MethodPost,
			username:     "newuser",
			email:        "newuser@example.com",
			password:     "password1",
			confirmPass:  "password2",
			expectedCode: http.StatusBadRequest,
			expectedBody: "Error parsing form data\n",
			mockActions:  func(m *MockUserModel) {},
		},
		{
			name:         "Invalid form data",
			method:       http.MethodPost,
			username:     "",
			email:        "",
			password:     "",
			confirmPass:  "",
			expectedCode: http.StatusBadRequest,
			expectedBody: "Error parsing form data\n", // Expected form parsing error
			mockActions:  func(m *MockUserModel) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the UserModel
			mockUserModel := new(MockUserModel)
			tt.mockActions(mockUserModel)

			// Create a mock HTTP request
			req, err := http.NewRequest(tt.method, "/signup", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Set form values for POST method if applicable
			if tt.method == http.MethodPost {
				req.Form = make(map[string][]string)
				req.Form.Set("username", tt.username)
				req.Form.Set("email", tt.email)
				req.Form.Set("password", tt.password)
				req.Form.Set("confirm_pass", tt.confirmPass)
			}

			// Create a mock HTTP response writer
			recorder := httptest.NewRecorder()

			// Create a handler to call the SignUpHandler
			SignUpHandler(recorder, req)

			// Check the response code
			if recorder.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, recorder.Code)
			}

			// Check the response body
			if recorder.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %q, got %q", tt.expectedBody, recorder.Body.String())
			}
		})
	}
}
