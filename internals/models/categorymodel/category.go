package categorymodel

import (
	"forum/internals/database"
)

type Category struct {
	ID   int64
	Name string
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

// Get all categories
func GetAllCategories() ([]Category, error) {
	query := `SELECT id, name FROM Categories`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
