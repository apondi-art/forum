package categorymodel

import (
	"database/sql"
	"fmt"

	"forum/internals/database"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var DefaultCategories = []Category{
	{ID: 1, Name: "Sports"},
	{ID: 2, Name: "Lifestyle"},
	{ID: 3, Name: "Education"},
	{ID: 4, Name: "Finance"},
	{ID: 5, Name: "Music"},
	{ID: 6, Name: "Culture"},
	{ID: 7, Name: "Technology"},
	{ID: 8, Name: "Health"},
	{ID: 9, Name: "Travel"},
	{ID: 10, Name: "Food"},
}

// GetAllCategories returns all categories from the database or defaults if none exist
func GetAllCategories() ([]Category, error) {
	query := `SELECT id, name FROM Categories ORDER BY name`
	rows, err := database.DB.Query(query)
	if err != nil {
		return DefaultCategories, nil
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, fmt.Errorf("failed to scan category: %v", err)
		}
		categories = append(categories, cat)
	}

	// If no categories in database, return defaults
	if len(categories) == 0 {
		return DefaultCategories, nil
	}

	return categories, nil
}

// CreateCategory creates a new category
func CreateCategory(name string) (int64, error) {
	query := `INSERT INTO Categories (name) VALUES (?)`
	result, err := database.DB.Exec(query, name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// AddCategoriesToPost adds multiple categories to a post in a transaction
func AddCategoriesToPost(tx *sql.Tx, postID int64, categoryIDs []int64) error {
	query := `INSERT INTO Post_Categories (post_id, category_id) VALUES (?, ?)`
	for _, categoryID := range categoryIDs {
		_, err := tx.Exec(query, postID, categoryID)
		if err != nil {
			return fmt.Errorf("failed to add category %d to post %d: %v", categoryID, postID, err)
		}
	}
	return nil
}
