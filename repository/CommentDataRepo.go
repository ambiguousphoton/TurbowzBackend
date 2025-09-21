package repository

import (
	"GoServer/models"
	"database/sql"
	"fmt"
)


type CommentRepo interface {
    CreateNewComment(vmd *models.CommentData) (int64,error)
	GetVideoComments(video_id int64) ([]models.CommentData, error)
}


type PostgresCommentRepo struct {
    db *sql.DB
}

func NewPostgresCommentRepo(db *sql.DB) CommentRepo {
    return &PostgresCommentRepo{db: db}
}


func (r *PostgresCommentRepo) CreateNewComment(new_comment *models.CommentData) (int64, error) {
	query := `
		INSERT INTO comments_table (parent_video_id, commenter_id, comment_text)
		VALUES ($1, $2, $3)
		RETURNING comment_id;
	`

	var comment_id int64
	err := r.db.QueryRow(query, new_comment.Parent_video_id, new_comment.Commenter_id, new_comment.Comment_text).
		Scan(&comment_id)
	
	if err != nil {
		return 0, fmt.Errorf("error Inserting the Comment %v", err)
	}


	return comment_id, nil
}


func (r *PostgresCommentRepo) GetVideoComments(video_id int64) ([]models.CommentData, error) {
	query := `
		SELECT 
			c.commenter_id,
			c.comment_text,
			c.created_at,
			u.user_handle
		FROM comments_table c
		JOIN user_data_table u ON c.commenter_id = u.user_id
		WHERE c.parent_video_id = $1
	`

	rows, err := r.db.Query(query, video_id)
	if err != nil {
		return nil, fmt.Errorf("error querying comments: %v", err)
	}
	defer rows.Close()

	var comments []models.CommentData
	for rows.Next() {
		var comment models.CommentData
		err := rows.Scan(&comment.Commenter_id, &comment.Comment_text, &comment.Comment_date, &comment.Commenter_Handle)
		if err != nil {
			return nil, fmt.Errorf("error scanning comment: %v", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over comments: %v", err)
	}

	return comments, nil
}