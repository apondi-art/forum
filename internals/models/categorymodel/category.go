package categorymodel

import (
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

// GetCategories returns all available categories
func GetAllCategories() ([]Category, error) {
	return DefaultCategories, nil
}

// Create a new category
func CreateCategory(name string) (int64, error) {
	query := `INSERT INTO Categories (name) VALUES (?)`
	result, err := database.DB.Exec(query, name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
