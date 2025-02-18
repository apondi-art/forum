package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"forum/internals/models/usermodel"

	"github.com/stretchr/testify/mock"
)

var DB *sql.DB

type MockUserModel struct {
	mock.Mock
}

func (m *MockUserModel) AuthenticateUser(email, password string) (bool, error) {
	args := m.Called(email, password)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserModel) GetUserByEmail(email string) (*usermodel.User, error) {
	args := m.Called(email)
	return args.Get(0).(*usermodel.User), args.Error(1)
}

func (m *MockUserModel) CreateSession(db *sql.DB, userID int64) (*usermodel.Session, error) {
	args := m.Called(db, userID)
	return args.Get(0).(*usermodel.Session), args.Error(1)
}

// Mocking a DB connection for tests
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Prepare(query string) (*sql.Stmt, error) {
	args := m.Called(query)
	return args.Get(0).(*sql.Stmt), args.Error(1)
}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		email        string
		password     string
		expectedCode int
		expectedBody string
		mockActions  func(m *MockUserModel)
	}{
		{
			name:         "InValid login",
			method:       http.MethodPost,
			email:        "user@example.com",
			password:     "validpassword",
			expectedCode: http.StatusBadRequest,
			expectedBody: "Failed to load error page\n",
			mockActions: func(m *MockUserModel) {
				m.On("AuthenticateUser", "user@example.com", "validpassword").Return(true, nil)
				m.On("GetUserByEmail", "user@example.com").Return(&usermodel.User{ID: 1}, nil)
				m.On("CreateSession", mock.Anything, int64(1)).Return(&usermodel.Session{ID: "session123"}, nil)
			},
		},
		{
			name:         "Missing form data",
			method:       http.MethodPost,
			email:        "",
			password:     "",
			expectedCode: http.StatusBadRequest,
			expectedBody: "Failed to load error page\n", // Should expect this message for missing data
			mockActions:  func(m *MockUserModel) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the UserModel
			mockUserModel := new(MockUserModel)
			tt.mockActions(mockUserModel)

			// Create a mock HTTP request
			req, err := http.NewRequest(tt.method, "/login", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Set form values for POST method if applicable
			if tt.method == http.MethodPost {
				req.Form = make(map[string][]string)
				req.Form.Set("email", tt.email)
				req.Form.Set("password", tt.password)
			}

			// Create a mock HTTP response writer
			recorder := httptest.NewRecorder()
			LoginHandler(recorder, req)
			if recorder.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, recorder.Code)
			}
			if recorder.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %q, got %q", tt.expectedBody, recorder.Body.String())
			}
		})
	}
}