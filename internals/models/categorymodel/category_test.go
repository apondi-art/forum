package categorymodel

import (
	"database/sql"
	"forum/internals/database"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func SetupTestDB(t *testing.T) (*sql.DB, func()) {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("Failed to open test database: %v", err)
    }

    // Create the Categories table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS Categories (
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL
        );
    `)
    if err != nil {
        t.Fatalf("Failed to create Categories table: %v", err)
    }

    // Create Post_Categories table for testing GetPostCategories
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS Post_Categories (
            post_id INTEGER,
            category_id INTEGER,
            PRIMARY KEY (post_id, category_id),
            FOREIGN KEY (category_id) REFERENCES Categories(id)
        );
    `)
    if err != nil {
        t.Fatalf("Failed to create Post_Categories table: %v", err)
    }

    // Restore global DB and return cleanup function
    oldDB := database.DB
    database.DB = db
    return db, func() {
        db.Close()
        database.DB = oldDB
    }
}


func TestSeedCategories(t *testing.T) {
    tests := []struct {
        name           string
        initialData    []Category
        expectedCount  int
        expectedError  bool
    }{
        {
            name:          "empty database gets seeded",
            initialData:   nil,
            expectedCount: len(DefaultCategories),
            expectedError: false,
        },
        {
            name:          "database with data doesn't get seeded",
            initialData:   []Category{{ID: 99, Name: "Test Category"}},
            expectedCount: 1, // Should keep existing data only
            expectedError: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            db, cleanup := SetupTestDB(t)
            defer cleanup()

            // Insert initial data if needed
            if len(tt.initialData) > 0 {
                for _, cat := range tt.initialData {
                    _, err := db.Exec("INSERT INTO Categories (id, name) VALUES (?, ?)", cat.ID, cat.Name)
                    if err != nil {
                        t.Fatalf("Failed to insert initial data: %v", err)
                    }
                }
            }

            // Call the function under test
            err := SeedCategories()

            // Check error expectations
            if tt.expectedError && err == nil {
                t.Errorf("Expected an error but got none")
            }
            if !tt.expectedError && err != nil {
                t.Errorf("Did not expect an error but got: %v", err)
            }

            // Verify the count of categories in DB
            var count int
            err = db.QueryRow("SELECT COUNT(*) FROM Categories").Scan(&count)
            if err != nil {
                t.Fatalf("Failed to count categories: %v", err)
            }
            
            if count != tt.expectedCount {
                t.Errorf("Expected %d categories in DB, but found %d", tt.expectedCount, count)
            }

            // If seeding should have happened, verify all default categories exist
            if tt.initialData == nil {
                for _, defaultCat := range DefaultCategories {
                    var name string
                    err = db.QueryRow("SELECT name FROM Categories WHERE id = ?", defaultCat.ID).Scan(&name)
                    if err != nil {
                        t.Errorf("Default category ID %d not found: %v", defaultCat.ID, err)
                    }
                    if name != defaultCat.Name {
                        t.Errorf("Default category name mismatch for ID %d: expected %s, got %s", 
                            defaultCat.ID, defaultCat.Name, name)
                    }
                }
            }
        })
    }
}

func TestGetPostCategories(t *testing.T) {
    tests := []struct {
        name             string
        postID           int64
        initialCategories []Category
        postCategoryLinks []struct{postID, categoryID int64}
        expectedResult   []string
        expectedError    bool
    }{
        {
            name:   "post with multiple categories",
            postID: 1,
            initialCategories: []Category{
                {ID: 1, Name: "Sports"},
                {ID: 2, Name: "Lifestyle"},
                {ID: 3, Name: "Education"},
            },
            postCategoryLinks: []struct{postID, categoryID int64}{
                {1, 1}, // Post 1 has Sports category
                {1, 3}, // Post 1 has Education category
            },
            expectedResult: []string{"Sports", "Education"},
            expectedError:  false,
        },
        {
            name:   "post with no categories",
            postID: 2,
            initialCategories: []Category{
                {ID: 1, Name: "Sports"},
                {ID: 2, Name: "Lifestyle"},
            },
            postCategoryLinks: []struct{postID, categoryID int64}{
                {1, 1}, // Only Post 1 has a category
            },
            expectedResult: []string{}, // Empty result
            expectedError:  false,
        },
        {
            name:   "non-existent post",
            postID: 999,
            initialCategories: []Category{
                {ID: 1, Name: "Sports"},
            },
            postCategoryLinks: []struct{postID, categoryID int64}{
                {1, 1}, // Only Post 1 has a category
            },
            expectedResult: []string{}, // Empty result for non-existent post
            expectedError:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            db, cleanup := SetupTestDB(t)
            defer cleanup()

            // Insert categories
            for _, cat := range tt.initialCategories {
                _, err := db.Exec("INSERT INTO Categories (id, name) VALUES (?, ?)", cat.ID, cat.Name)
                if err != nil {
                    t.Fatalf("Failed to insert categories: %v", err)
                }
            }

            // Insert post-category links
            for _, link := range tt.postCategoryLinks {
                _, err := db.Exec("INSERT INTO Post_Categories (post_id, category_id) VALUES (?, ?)", 
                    link.postID, link.categoryID)
                if err != nil {
                    t.Fatalf("Failed to insert post-category links: %v", err)
                }
            }

            // Call the function under test
            result, err := GetPostCategories(tt.postID)

            // Check error expectations
            if tt.expectedError && err == nil {
                t.Errorf("Expected an error but got none")
            }
            if !tt.expectedError && err != nil {
                t.Errorf("Did not expect an error but got: %v", err)
            }

            // Sort both slices to ensure consistent comparison
            // (Assuming order doesn't matter for category names)
            // Note: You might need to implement a function to sort strings if not imported
            
            // Compare lengths first
            if len(result) != len(tt.expectedResult) {
                t.Errorf("Expected %d categories but got %d", len(tt.expectedResult), len(result))
            } else {
                // Simple check for each expected item
                for _, expected := range tt.expectedResult {
                    found := false
                    for _, actual := range result {
                        if expected == actual {
                            found = true
                            break
                        }
                    }
                    if !found {
                        t.Errorf("Expected category '%s' not found in result", expected)
                    }
                }
            }
        })
    }
}