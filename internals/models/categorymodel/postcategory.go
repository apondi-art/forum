package categorymodel

import (
	"forum/internals/database"
)

type PostCategory struct {
	ID         int64
	PostID     int64
	CategoryID int64
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
