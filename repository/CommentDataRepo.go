package repository

import (
	"GoServer/models"
	"database/sql"
	"fmt"
)


type CommentRepo interface {
    CreateNewComment(vmd *models.CommentData) (int64,error)
	GetVideoComments(video_id int64, limit, offset int) ([]models.CommentData, error)
	GetVideoCommentsCount(video_id int64) (int64, error)
	CreateNewEcoComment(ecd *models.EcoCommentData) (int64, error)
	GetEcoComments(eco_id int64, limit, offset int) ([]models.EcoCommentData, error)
	GetEcoCommentsCount(eco_id int64) (int64, error)
}


type PostgresCommentRepo struct {
    db *sql.DB
}

func NewPostgresCommentRepo(db *sql.DB) CommentRepo {
    return &PostgresCommentRepo{db: db}
}


func (r *PostgresCommentRepo) CreateNewComment(new_comment *models.CommentData) (int64, error) {
	query := `
		INSERT INTO comments_table (parent_video_id, commenter_id, comment_text, parent_comment_id)
		VALUES ($1, $2, $3, $4)
		RETURNING comment_id;
	`

	var comment_id int64
	err := r.db.QueryRow(query, new_comment.Parent_video_id, new_comment.Commenter_id, new_comment.Comment_text, new_comment.Parent_Comment_ID).
		Scan(&comment_id)
	
	if err != nil {
		return 0, fmt.Errorf("error Inserting the Comment %v", err)
	}

	fmt.Printf("Comment created - ID: %d, VideoID: %d, ParentCommentID: %v\n", comment_id, new_comment.Parent_video_id, new_comment.Parent_Comment_ID)
	return comment_id, nil
}


func (r *PostgresCommentRepo) GetVideoComments(video_id int64, limit, offset int) ([]models.CommentData, error) {
	query := `
		SELECT 
			c.commenter_id,
			c.comment_text,
			c.created_at,
			u.user_handle,
			u.user_profile_name,
			c.comment_id,
			c.parent_comment_id
		FROM comments_table c
		JOIN user_data_table u ON c.commenter_id = u.user_id
		WHERE c.parent_video_id = $1
		ORDER BY c.parent_comment_id NULLS FIRST, c.created_at
		LIMIT $2 OFFSET $3;

	`

	rows, err := r.db.Query(query, video_id, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying comments: %v", err)
	}
	defer rows.Close()

	var comments []models.CommentData
	for rows.Next() {
		var comment models.CommentData
		err := rows.Scan(&comment.Commenter_id, &comment.Comment_text, &comment.Comment_date, &comment.Commenter_Handle, &comment.Commenter_Name, &comment.Comment_id, &comment.Parent_Comment_ID)
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

func (r *PostgresCommentRepo) CreateNewEcoComment(new_comment *models.EcoCommentData) (int64, error) {
	query := `
		INSERT INTO eco_comments_table (parent_eco_id, commenter_id, comment_text)
		VALUES ($1, $2, $3)
		RETURNING comment_id;
	`

	var comment_id int64
	err := r.db.QueryRow(query, new_comment.Parent_Eco_id, new_comment.Commenter_id, new_comment.Comment_text).
		Scan(&comment_id)
	
	if err != nil {
		return 0, fmt.Errorf("error inserting eco comment %v", err)
	}

	return comment_id, nil
}

func (r *PostgresCommentRepo) GetEcoComments(eco_id int64, limit, offset int) ([]models.EcoCommentData, error) {
	query := `
		SELECT 
			c.commenter_id,
			c.comment_text,
			c.created_at,
			u.user_handle,
			u.user_profile_name
		FROM eco_comments_table c
		JOIN user_data_table u ON c.commenter_id = u.user_id
		WHERE c.parent_eco_id = $1
		ORDER BY c.created_at
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.db.Query(query, eco_id, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying eco comments: %v", err)
	}
	defer rows.Close()

	var comments []models.EcoCommentData
	for rows.Next() {
		var comment models.EcoCommentData
		err := rows.Scan(&comment.Commenter_id, &comment.Comment_text, &comment.Comment_date, &comment.Commenter_Handle, &comment.Commenter_Name)
		if err != nil {
			return nil, fmt.Errorf("error scanning eco comment: %v", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over eco comments: %v", err)
	}

	return comments, nil
}

func (r *PostgresCommentRepo) GetVideoCommentsCount(video_id int64) (int64, error) {
	query := `SELECT COUNT(*) FROM comments_table WHERE parent_video_id = $1`
	var count int64
	err := r.db.QueryRow(query, video_id).Scan(&count)
	return count, err
}

func (r *PostgresCommentRepo) GetEcoCommentsCount(eco_id int64) (int64, error) {
	query := `SELECT COUNT(*) FROM eco_comments_table WHERE parent_eco_id = $1`
	var count int64
	err := r.db.QueryRow(query, eco_id).Scan(&count)
	return count, err
}