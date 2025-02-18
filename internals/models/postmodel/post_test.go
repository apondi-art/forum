package postmodel

import (
	"errors"
	"testing"

	"forum/internals/database"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreatePost(t *testing.T) {
	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer mockDB.Close()

	// Replace the real database with our mock
	originalDB := database.DB
	database.DB = mockDB
	defer func() { database.DB = originalDB }()

	// Test data
	userID := int64(123)
	title := "Test Post"
	content := "This is test content"
	categoryIDs := []int64{1, 2, 3}
	expectedPostID := int64(456)

	// Set up the mock expectations
	mock.ExpectBegin()

	// Expect the INSERT INTO Posts query
	mock.ExpectExec("INSERT INTO Posts").
		WithArgs(userID, title, content).
		WillReturnResult(sqlmock.NewResult(expectedPostID, 1))

	// Expect the INSERT INTO Post_Categories queries for each category
	for _, categoryID := range categoryIDs {
		mock.ExpectExec("INSERT INTO Post_Categories").
			WithArgs(expectedPostID, categoryID).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	// Expect a transaction commit
	mock.ExpectCommit()

	// Call the function being tested
	postID, err := CreatePost(userID, title, content, categoryIDs)
	// Check for errors
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check the returned post ID
	if postID != expectedPostID {
		t.Errorf("Expected post ID %d, but got %d", expectedPostID, postID)
	}

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// Additional test for error handling during post insertion
func TestCreatePostInsertError(t *testing.T) {
	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer mockDB.Close()

	// Replace the real database with our mock
	originalDB := database.DB
	database.DB = mockDB
	defer func() { database.DB = originalDB }()

	// Test data
	userID := int64(123)
	title := "Test Post"
	content := "This is test content"
	categoryIDs := []int64{1, 2, 3}
	expectedError := errors.New("database error")

	// Set up the mock expectations
	mock.ExpectBegin()

	// Expect the INSERT INTO Posts query to fail
	mock.ExpectExec("INSERT INTO Posts").
		WithArgs(userID, title, content).
		WillReturnError(expectedError)

	// Expect a rollback since there was an error
	mock.ExpectRollback()

	// Call the function being tested
	_, err = CreatePost(userID, title, content, categoryIDs)

	// Check for errors
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// Additional test for error handling during category insertion
func TestCreatePostCategoryInsertError(t *testing.T) {
	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer mockDB.Close()

	// Replace the real database with our mock
	originalDB := database.DB
	database.DB = mockDB
	defer func() { database.DB = originalDB }()

	// Test data
	userID := int64(123)
	title := "Test Post"
	content := "This is test content"
	categoryIDs := []int64{1, 2, 3}
	expectedPostID := int64(456)
	expectedError := errors.New("category insert error")

	// Set up the mock expectations
	mock.ExpectBegin()

	// Expect the INSERT INTO Posts query
	mock.ExpectExec("INSERT INTO Posts").
		WithArgs(userID, title, content).
		WillReturnResult(sqlmock.NewResult(expectedPostID, 1))

	// First category insert succeeds
	mock.ExpectExec("INSERT INTO Post_Categories").
		WithArgs(expectedPostID, categoryIDs[0]).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Second category insert fails
	mock.ExpectExec("INSERT INTO Post_Categories").
		WithArgs(expectedPostID, categoryIDs[1]).
		WillReturnError(expectedError)

	// Expect a rollback since there was an error
	mock.ExpectRollback()

	// Call the function being tested
	_, err = CreatePost(userID, title, content, categoryIDs)

	// Check for errors
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
