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
			name:         "InValid signup",
			method:       http.MethodPost,
			username:     "newuser",
			email:        "newuser@example.com",
			password:     "invalidpassword",
			confirmPass:  "invalidpassword",
			expectedCode: http.StatusBadRequest, 
			expectedBody: "Failed to load error page\n",           
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
			expectedBody: "Failed to load error page\n",
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
			expectedBody: "Failed to load error page\n", // Expected form parsing error
			mockActions:  func(m *MockUserModel) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the UserModel
			mockUserModel := new(MockUserModel)
			tt.mockActions(mockUserModel)
			req, err := http.NewRequest(tt.method, "/signup", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			// Create a mock HTTP response writer
			recorder := httptest.NewRecorder()
			SignUpHandler(recorder, req)
			if recorder.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, recorder.Code)
			}
			if recorder.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %q, got %q", tt.expectedBody, recorder.Body.String())
			}
		})
	}
}
