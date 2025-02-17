//This code is a unit test for the GetAllCategories function using mocking with the sqlmock library. 
//It tests how the function behaves under different scenarios, such as successful queries, empty results, 
//database errors, and scan errors. Let’s break it down step by step:


/*    database/sql: For interacting with SQL databases.
    fmt: For formatting strings.
    reflect: For comparing complex data structures (like slices).
    testing: Go’s built-in testing framework.
    github.com/DATA-DOG/go-sqlmock: A library for mocking SQL databases.*/
package categorymodel

import (
	"fmt"
	"reflect"
	"testing"

	"forum/internals/database"

	"github.com/DATA-DOG/go-sqlmock"
)
/*    name: A descriptive name for the test case.
    mockSetup: A function that sets up the mock database behavior for this test case.
    expectedResult: The expected result from GetAllCategories.
    expectedError: Whether an error is expected from GetAllCategories.*/

func TestGetAllCategories(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		mockSetup      func(mock sqlmock.Sqlmock)
		expectedResult []Category
		expectedError  bool
	}{
		{
			/*    Scenario: The database query returns 3 categories.

    Mock Setup:
        Create mock rows with 3 categories.
        Expect the query SELECT id, name FROM Categories ORDER BY name and return the mock rows.
    Expected Result: The function should return the 3 categories.
    Expected Error: No error is expected.*/
			name: "successful query with categories",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Sports").
					AddRow(2, "Lifestyle").
					AddRow(3, "Education")
				mock.ExpectQuery("SELECT id, name FROM Categories ORDER BY name").
					WillReturnRows(rows)
			},
			expectedResult: []Category{
				{ID: 1, Name: "Sports"},
				{ID: 2, Name: "Lifestyle"},
				{ID: 3, Name: "Education"},
			},
			expectedError: false,
		},
		{
			name: "empty result returns default categories",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"})
				mock.ExpectQuery("SELECT id, name FROM Categories ORDER BY name").
					WillReturnRows(rows)
			},
			expectedResult: DefaultCategories,
			expectedError:  false,
		},
		{
			name: "database error returns default categories",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name FROM Categories ORDER BY name").
					WillReturnError(fmt.Errorf("database connection error"))
			},
			expectedResult: DefaultCategories,
			expectedError:  false,
		},
		{
			name: "scan error returns error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("invalid", "Sports") // This will cause a scan error
				mock.ExpectQuery("SELECT id, name FROM Categories ORDER BY name").
					WillReturnRows(rows)
			},
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock database connection
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database: %v", err)
			}
			defer db.Close()

			// Replace the global DB with our mock
			oldDB := database.DB
			database.DB = db
			defer func() { database.DB = oldDB }()

			// Set up the mock expectations
			tt.mockSetup(mock)

			// Call the function under test
			result, err := GetAllCategories()

			// Check error expectations
			if tt.expectedError && err == nil {
				t.Error("Expected an error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Did not expect an error but got: %v", err)
			}

			// For error cases, we don't need to check the result
			if tt.expectedError {
				return
			}

			// Check that result matches expected
			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("Expected %+v but got %+v", tt.expectedResult, result)
			}

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
