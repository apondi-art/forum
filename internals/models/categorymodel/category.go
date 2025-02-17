package categorymodel

import (
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

// SeedCategories inserts default categories if they don't exist
func SeedCategories() error {
	var count int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM Categories`).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check categories count: %v", err)
	}

	if count == 0 { // If no categories exist, insert defaults
		tx, err := database.DB.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %v", err)
		}
		defer tx.Rollback()

		query := `INSERT INTO Categories (id, name) VALUES (?, ?)`
		for _, category := range DefaultCategories {
			_, err := tx.Exec(query, category.ID, category.Name)
			if err != nil {
				return fmt.Errorf("failed to insert category %s: %v", category.Name, err)
			}
		}

		if err = tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %v", err)
		}
	}
	return nil
}

// GetPostCategories retrieves category names for a post
func GetPostCategories(postID int64) ([]string, error) {
	query := `
        SELECT c.name
        FROM Categories c
        JOIN Post_Categories pc ON c.id = pc.category_id
        WHERE pc.post_id = ?
    `
	rows, err := database.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
