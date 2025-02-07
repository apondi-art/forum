package postmodel

import "forum/internals/database"

// GetAll retrieves all posts from the database, ordered by creation date in descending order.
// It executes a SQL query to fetch all posts and scans the results into a slice of Post structs.
// Returns the slice of posts and any error encountered during database operations.
// If no posts are found, returns an empty slice with nil error.
func GetAll() ([]Post, error) {
	query := `
        SELECT id, user_id, title, content, created_at, updated_at
        FROM Posts
        ORDER BY created_at DESC
    `
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content,
			&post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// GetUserPosts retrieves all posts created by a specific user from the database.
func GetUserPosts(userID int64) ([]Post, error) {
	query := `
        SELECT id, user_id, title, content, created_at, updated_at
        FROM Posts
        WHERE user_id = ?
        ORDER BY created_at DESC
    `
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title,
			&post.Content, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostsByCategory(categoryID int64) ([]Post, error) {
	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.updated_at
        FROM Posts p
        JOIN Post_Categories pc ON p.id = pc.post_id
        WHERE pc.category_id = ?
        ORDER BY p.created_at DESC
    `
	rows, err := database.DB.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title,
			&post.Content, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
