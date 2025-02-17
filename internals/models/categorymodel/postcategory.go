package categorymodel

import (
	"forum/internals/database"
	"forum/internals/models/usermodel"
	"forum/internals/models/viewmodel"
)

// GetPostsBySingleCategory retrieves posts for a single category or all posts if categoryID is 0
// If showLiked is true, it only returns posts liked by the specified user
func GetPostsBySingleCategory(categoryID int64, userID int64, showLiked bool) ([]viewmodel.PostView, error) {
	var posts []viewmodel.PostView

	// Base query for all posts or filtered by category
	query := `
        SELECT DISTINCT p.id, p.title, p.content, u.username, p.created_at, p.updated_at
        FROM Posts p
        JOIN Users u ON p.user_id = u.id
    `

	// Add category filter if categoryID is not 0
	var args []interface{}

	// Add liked posts filter if requested
	if showLiked {
		query += `
            JOIN likes_dislikes ld ON p.id = ld.post_id
            WHERE ld.user_id = ? AND ld.reaction_type = 'like'
        `
		args = append(args, userID)

		// If category is also selected, add AND condition
		if categoryID != 0 {
			query += `
                AND p.id IN (
                    SELECT post_id FROM Post_Categories 
                    WHERE category_id = ?
                )
            `
			args = append(args, categoryID)
		}
	} else if categoryID != 0 {
		// If only category filter is active
		query += `
            JOIN Post_Categories pc ON p.id = pc.post_id
            WHERE pc.category_id = ?
        `
		args = append(args, categoryID)
	}

	query += " ORDER BY p.created_at DESC"

	// Execute query
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process results
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

		// Get categories for the post
		categories, err := GetPostCategories(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		// Get reaction counts
		likeCount, dislikeCount, err := usermodel.GetReactionCounts(post.ID)
		if err != nil {
			return nil, err
		}
		post.LikeCount = likeCount
		post.DislikeCount = dislikeCount

		// Get comments
		comments, err := viewmodel.GetPostComments(post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments

		posts = append(posts, post)
	}

	return posts, nil
}
