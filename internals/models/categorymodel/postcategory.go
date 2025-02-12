package categorymodel

import (
	"fmt"
	"strings"

	"forum/internals/database"
)

type PostCategory struct {
	ID         int64
	PostID     int64
	CategoryID int64
}

// Add a category to a post
func AddCategoryToPost(postID, categoryID int64) error {
	query := `INSERT INTO Post_Categories (post_id, category_id) VALUES (?, ?)`
	_, err := database.DB.Exec(query, postID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to add category to post: %v", err)
	}
	return nil
}

// Remove a category from a post
func RemoveCategoryFromPost(postID, categoryID int64) error {
	query := `DELETE FROM Post_Categories WHERE post_id = ? AND category_id = ?`
	_, err := database.DB.Exec(query, postID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to remove category from post: %v", err)
	}
	return nil
}

// Get all posts for a category
func GetCategoryPosts(categoryID int64) ([]int64, error) {
	query := `
        SELECT post_id
        FROM Post_Categories
        WHERE category_id = ?
    `
	rows, err := database.DB.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postIDs []int64
	for rows.Next() {
		var postID int64
		if err := rows.Scan(&postID); err != nil {
			return nil, err
		}
		postIDs = append(postIDs, postID)
	}
	return postIDs, nil
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

// GetPostsByCategories retrieves posts that belong to the specified categories
func GetPostsByCategories(categoryIDs []int64) ([]int64, error) {
	if len(categoryIDs) == 0 {
		return nil, fmt.Errorf("no category IDs provided")
	}

	placeholders := strings.Repeat("?,", len(categoryIDs))
	placeholders = placeholders[:len(placeholders)-1] // Remove trailing comma

	query := fmt.Sprintf(`
        SELECT DISTINCT post_id
        FROM Post_Categories
        WHERE category_id IN (%s)
    `, placeholders)

	args := make([]interface{}, len(categoryIDs))
	for i, id := range categoryIDs {
		args[i] = id
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postIDs []int64
	for rows.Next() {
		var postID int64
		if err := rows.Scan(&postID); err != nil {
			return nil, err
		}
		postIDs = append(postIDs, postID)
	}
	return postIDs, nil
}
