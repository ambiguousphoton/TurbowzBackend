package repository

import (
	"GoServer/models"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type EventRepo interface {
	CreateNewEvent(event *models.EventMetaData) (int64, error)
	GetEventByID(event_id int64) (*models.EventMetaData, error)
	GetEventEngagementInfo(event_id int64) (int64, int64, int64, int64, float64, string, error)
	SaveEventTrendingDelta(event_id int64, currentScore float64, delta float64) error
	GetAllEventIDs(limit int, offset int) ([]int64, error)
	GetTrendingEventIDs(limit int, offset int)([]int64, error)
	GetTrendingEvents(limit int, offset int)([]models.EventMetaData, error)
	IncrementViewsOfEvent(event_id int64)( error)
}

type PostgresEventRepo struct {
	db *sql.DB
}

func NewPostgresEventRepo(db *sql.DB) EventRepo {
	return &PostgresEventRepo{db: db}
}

func (r *PostgresEventRepo) CreateNewEvent(event *models.EventMetaData) (int64, error) {
	query := `
		INSERT INTO event_data (event_url, event_title, uploader_id, event_description, tags, images_count, event_start_time, event_end_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING event_id;
	`

	var event_id int64
	err := r.db.QueryRow(query, event.Event_Url, event.Event_Title, event.Uploader_ID, event.Event_Description, pq.Array(event.Tags), event.Images_Count, event.Event_Start_Time, event.Event_End_Time).
		Scan(&event_id)

	if err != nil {
		return 0, fmt.Errorf("error inserting event: %v", err)
	}

	fmt.Printf("Event created - ID: %d, UploaderID: %d\n", event_id, event.Uploader_ID)
	return event_id, nil
}

func (r *PostgresEventRepo) GetEventByID(event_id int64) (*models.EventMetaData, error) {
	_, err := r.db.Exec(`UPDATE event_data SET view_count = view_count + 1 WHERE event_id = $1`, event_id)
	if err != nil {
		return nil, fmt.Errorf("error updating view count: %v", err)
	}

	query := `
		SELECT 
			e.event_id,
			e.event_url,
			e.event_title,
			e.uploader_id,
			u.user_handle,
			u.user_profile_name,
			e.event_description,
			e.view_count,
			e.luv_count,
			e.comment_count,
			e.saves_count,
			e.images_count,
			e.tags,
			e.created_at,
			e.event_start_time,
			e.event_end_time
		FROM event_data e
		JOIN user_data_table u ON e.uploader_id = u.user_id
		WHERE e.event_id = $1
	`

	var event models.EventMetaData
	err = r.db.QueryRow(query, event_id).Scan(
		&event.Event_Id,
		&event.Event_Url,
		&event.Event_Title,
		&event.Uploader_ID,
		&event.Uploader_Handle,
		&event.Uploader_Name,
		&event.Event_Description,
		&event.View_Count,
		&event.Luv_Count,
		&event.Comment_Count,
		&event.Saves_Count,
		&event.Images_Count,
		pq.Array(&event.Tags),
		&event.Created_At,
		&event.Event_Start_Time,
		&event.Event_End_Time,
	)

	if err != nil {
		return nil, fmt.Errorf("error querying event: %v", err)
	}

	return &event, nil
}

func (r *PostgresEventRepo) GetEventEngagementInfo(event_id int64) (int64, int64, int64, int64, float64, string, error) {
	query := `SELECT view_count, luv_count, comment_count, saves_count, trending_score, created_at FROM event_data WHERE event_id = $1`
	var views, likes, comments, shares int64
	var lastScore float64
	var timestamp string
	err := r.db.QueryRow(query, event_id).Scan(&views, &likes, &comments, &shares, &lastScore, &timestamp)
	return views, likes, comments, shares, lastScore, timestamp, err
}

func (r *PostgresEventRepo) SaveEventTrendingDelta(event_id int64, currentScore float64, delta float64) error {
	query := `UPDATE event_data SET trending_score = $1, trending_delta = $2 WHERE event_id = $3`
	_, err := r.db.Exec(query, currentScore, delta, event_id)
	return err
}

func (r *PostgresEventRepo) GetAllEventIDs(limit int, offset int) ([]int64, error) {
	query := `SELECT event_id FROM event_data ORDER BY event_id LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}




func (r *PostgresEventRepo) GetTrendingEventIDs(limit int, offset int) ([]int64, error) {
	query := `SELECT event_id FROM event_data ORDER BY trending_delta DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *PostgresEventRepo) GetTrendingEvents(limit int, offset int) ([]models.EventMetaData, error) {
	query := `
		SELECT 
			e.event_id,
			e.event_url,
			e.event_title,
			e.uploader_id,
			u.user_handle,
			u.user_profile_name,
			e.event_description,
			e.view_count,
			e.luv_count,
			e.comment_count,
			e.saves_count,
			e.images_count,
			e.tags,
			e.created_at,
			e.event_start_time,
			e.event_end_time
		FROM event_data e
		JOIN user_data_table u ON e.uploader_id = u.user_id
		ORDER BY e.trending_delta DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.EventMetaData
	for rows.Next() {
		var event models.EventMetaData
		err := rows.Scan(
			&event.Event_Id,
			&event.Event_Url,
			&event.Event_Title,
			&event.Uploader_ID,
			&event.Uploader_Handle,
			&event.Uploader_Name,
			&event.Event_Description,
			&event.View_Count,
			&event.Luv_Count,
			&event.Comment_Count,
			&event.Saves_Count,
			&event.Images_Count,
			pq.Array(&event.Tags),
			&event.Created_At,
			&event.Event_Start_Time,
			&event.Event_End_Time,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}


func (r *PostgresEventRepo) IncrementViewsOfEvent(eventID int64) error {
	query := `
		UPDATE events
		SET view_count = view_count + 1
		WHERE event_id = $1
	`

	result, err := r.db.Exec(query, eventID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event not found: %d", eventID)
	}

	return nil
}
