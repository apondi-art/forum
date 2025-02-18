package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
    
	"forum/internals/models/categorymodel"
	"forum/internals/models/viewmodel"
)

type  NewHomepage struct {
	auth     AuthInterface 
	category CategoryModel 
	post     PostModel     
}


type AuthInterface interface {
	GetUserFromSession(r *http.Request) (int64, bool)
	GetUserNameByID(id int64) (string, error)
}

type CategoryModel interface {
	GetAllCategories() ([]categorymodel.Category, error)
	GetPostsByCategories(categoryIDs []int64) ([]int64, error)
}

type PostModel interface {
	GetFilteredPosts(userID int64, categoryID sql.NullInt64, isPrivate bool) ([]viewmodel.PostView, error)
	GetPostsByIDs(postIDs []int64) ([]viewmodel.PostView, error)
}

// Implement the ServeHTTP method for your handler
func (h  NewHomepage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
_, isLoggedIn := h.auth.GetUserFromSession(r)

if !isLoggedIn { //Guest user
	 http.Error(w, "Guest User", http.StatusOK) 
} else { //Logged in user
	 http.Error(w, "Logged in user", http.StatusOK)
}

}
// Mock implementations
type mockAuth struct {
	userID     int64
	isLoggedIn bool
	userName   string
}

func (m *mockAuth) GetUserFromSession(r *http.Request) (int64, bool) {
	return m.userID, m.isLoggedIn
}

func (m *mockAuth) GetUserNameByID(id int64) (string, error) {
	return m.userName, nil
}

type mockCategoryModel struct {
	categories []categorymodel.Category
	postIDs    []int64
	err        error
}

func (m *mockCategoryModel) GetAllCategories() ([]categorymodel.Category, error) {
	return m.categories, m.err
}

func (m *mockCategoryModel) GetPostsByCategories(categoryIDs []int64) ([]int64, error) {
	return m.postIDs, m.err
}

type mockPostModel struct {
	posts []viewmodel.PostView
	err   error
}

func (m *mockPostModel) GetFilteredPosts(userID int64, categoryID sql.NullInt64, isPrivate bool) ([]viewmodel.PostView, error) {
	return m.posts, m.err
}

func (m *mockPostModel) GetPostsByIDs(postIDs []int64) ([]viewmodel.PostView, error) {
	return m.posts, m.err
}

func TestHomepage(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		mockAuth     mockAuth
		mockCategory mockCategoryModel
		mockPost     mockPostModel
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Invalid Path Returns 404",
			path:         "/invalid",
			expectedCode: http.StatusNotFound,
		},
		{
			name: "Successful Load For Logged In User",
			path: "/",
			mockAuth: mockAuth{
				userID:     1,
				isLoggedIn: true,
				userName:   "testUser",
			},
			mockCategory: mockCategoryModel{
				categories: []categorymodel.Category{
					{ID: 1, Name: "Test Category"},
				},
				err: nil,
			},
			mockPost: mockPostModel{
				posts: []viewmodel.PostView{
					{ID: 1, Title: "Test Post"},
				},
				err: nil,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Successful Load For Guest User",
			path: "/",
			mockAuth: mockAuth{
				isLoggedIn: false,
			},
			mockCategory: mockCategoryModel{
				categories: []categorymodel.Category{
					{ID: 1, Name: "Test Category"},
				},
				err: nil,
			},
			mockPost: mockPostModel{
				posts: []viewmodel.PostView{
					{ID: 1, Title: "Test Post"},
				},
				err: nil,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Category Load Error",
			path: "/",
			mockCategory: mockCategoryModel{
				err: sql.ErrConnDone,
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Guest User\n",
		},
		{
			name: "Guest User",
			path: "/",
			mockCategory: mockCategoryModel{
				categories: []categorymodel.Category{
					{ID: 1, Name: "Test Category"},
				},
				err: nil,
			},
			mockPost: mockPostModel{
				err: sql.ErrConnDone,
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Guest User\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup request
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			handler := NewHomepage{
				auth:     &tt.mockAuth,
				category: &tt.mockCategory,
				post:     &tt.mockPost,
			}
			// Call handler
			handler.ServeHTTP(w, req)

			if tt.expectedBody != "" && w.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}
