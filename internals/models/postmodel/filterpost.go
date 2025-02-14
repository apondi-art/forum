package postmodel

import (
	"database/sql"

	"forum/internals/database"
	"forum/internals/models/categorymodel"
	"forum/internals/models/usermodel"
	"forum/internals/models/viewmodel"
)

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

// GetFilteredPosts retrieves posts based on various filters
func GetFilteredPosts(userID int64, categoryID sql.NullInt64, showLiked bool) ([]viewmodel.PostView, error) {
	query := `
        SELECT DISTINCT p.id, p.title, p.content, u.username, p.created_at, p.updated_at
        FROM Posts p
        JOIN Users u ON p.user_id = u.id
    `
	var params []interface{}

	// Add category filter if specified
	if categoryID.Valid {
		query += `
            JOIN Post_Categories pc ON p.id = pc.post_id
            WHERE pc.category_id = ?
        `
		params = append(params, categoryID.Int64)
	}

	// Add liked posts filter if specified
	if showLiked {
		if categoryID.Valid {
			query += " AND "
		} else {
			query += " WHERE "
		}
		query += `
            p.id IN (
                SELECT post_id 
                FROM Likes_Dislikes 
                WHERE user_id = ? AND reaction_type = 'like'
            )
        `
		params = append(params, userID)
	}

	query += " ORDER BY p.created_at DESC"

	rows, err := database.DB.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []viewmodel.PostView
	for rows.Next() {
		var post viewmodel.PostView
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Author,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get additional data for each post
		post.LikeCount, post.DislikeCount, err = usermodel.GetReactionCounts(post.ID)
		if err != nil {
			return nil, err
		}

		// Get categories for the post
		post.Categories, err = categorymodel.GetPostCategories(post.ID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
