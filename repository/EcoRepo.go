package repository

import (
        "database/sql"
        "log"
        "GoServer/models"

		"fmt"
		
		"github.com/lib/pq"
        )

type EcoRepo interface {
	CreateEcoPost(eco *models.EcoMetaData)(int64, error)
	SearchEcosByUserID(userID int64)([]models.EcoMetaData, error)
    RecommendEcosByUserEmbedding(userID int64, limit, offset int) ([]models.EcoMetaData, error)
	UpdateLuv(ecoID, userID int64) (luvved bool, err error)
	LuvStatus(ecoID, userID int64) (bool, int64, error)
	GetEcoMetaData(ecoID int64) (models.EcoMetaData, error)
	GetAllEcoIDs(limit, offset int) ([]int64, error)
	GetEcoEngagementInfo(ecoID int64) (views, luvs, comments, shares int64, lastScore float64, timestampStr string, err error)
	SaveEcoTrendingDelta(ecoID int64, score float64, delta float64) error
	GetTrendingEcos(limit, offset int) ([]models.EcoMetaData, error)	
	GetEchoScore(echoID int64) (models.EchoScore, error)
}

type PostgresEcoRepo struct {
    db *sql.DB
}

func NewPostgresEcoRepo(db *sql.DB) EcoRepo {
    return &PostgresEcoRepo{db: db}
}

func (r *PostgresEcoRepo) CreateEcoPost(eco *models.EcoMetaData) (int64, error) {
	query := `
		INSERT INTO eco_data (
			uploader_id,
			eco_text,
			tags,
			images_count,
			eco_url
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING eco_id;
	`

	var ecoID int64


	err := r.db.QueryRow(
		query,
		eco.Uploader_ID,
		eco.Eco_Text,
		pq.Array(eco.Tags),
		eco.Images_Count,
		eco.Eco_Url,
	).Scan(&ecoID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert eco post: %w", err)
	}

	log.Printf("Eco post created with ID: %d", ecoID)
	return ecoID, nil
}

func (r *PostgresEcoRepo) SearchEcosByUserID(userID int64)([]models.EcoMetaData, error) {
    query := `
                SELECT 
                e.eco_id,
                e.eco_text,
                e.eco_url,
                e.images_count,
                e.created_at,
				e.view_count,
				e.tags,
				e.luv_count,
				e.comment_count,
                u.user_profile_name AS uploader_name,
                u.user_handle AS uploader_handle,
				u.user_id
                FROM eco_data e
                JOIN user_data_table u
                ON e.uploader_id = u.user_id
                WHERE e.uploader_id = $1;
    `
    rows, err := r.db.Query(query, userID)
    if err != nil {
        log.Printf("SearchEcosByUserID: Failed to execute search query for userID '%d' - %v", userID, err)
        return nil, err
    }
    defer rows.Close()  
    var results []models.EcoMetaData;
    for rows.Next(){
        var ecoData models.EcoMetaData
        if err := rows.Scan(
            &ecoData.Eco_Id,
            &ecoData.Eco_Text,
            &ecoData.Eco_Url,
            &ecoData.Images_Count,
            &ecoData.Created_At,
			&ecoData.View_Count,
			pq.Array(&ecoData.Tags),
			&ecoData.Luv_Count,
			&ecoData.Comment_Count,
            &ecoData.Uploader_Name,
            &ecoData.Uploader_Handle,
			&ecoData.Uploader_ID,
        ); err != nil {
            log.Printf("SearchEcosByUserID: Failed to scan row for userID '%d' - %v", userID, err)
            return nil, err
        }
		ecoData.Uploader_ID = userID
        results = append(results, ecoData)
    }
    return results, nil
}

func (r *PostgresEcoRepo) RecommendEcosByUserEmbedding(userID int64, limit, offset int) ([]models.EcoMetaData, error) {
	query := `
		SELECT 
			e.eco_id,
			e.eco_text,
			e.eco_url,
			e.images_count,
			e.created_at,
			e.view_count,
			e.tags,
			e.luv_count,
			e.comment_count,
			u.user_profile_name AS uploader_name,
			u.user_handle AS uploader_handle,
			u.user_id
		FROM eco_data e
		JOIN user_data_table u ON e.uploader_id = u.user_id
		JOIN user_data_table usr ON usr.user_id = $1
		ORDER BY usr.eco_embeddings <=> e.embeddings ASC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		log.Printf("RecommendEcosByUserEmbedding: query failed for userID '%d' - %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var results []models.EcoMetaData
	for rows.Next() {
		var ecoData models.EcoMetaData
		err := rows.Scan(
			&ecoData.Eco_Id,
			&ecoData.Eco_Text,
			&ecoData.Eco_Url,
			&ecoData.Images_Count,
			&ecoData.Created_At,
			&ecoData.View_Count,
			pq.Array(&ecoData.Tags),
			&ecoData.Luv_Count,
			&ecoData.Comment_Count,
			&ecoData.Uploader_Name,
			&ecoData.Uploader_Handle,
			&ecoData.Uploader_ID,
		)
		if err != nil {
			log.Printf("RecommendEcosByUserEmbedding: row scan failed for userID '%d' - %v", userID, err)
			return nil, err
		}
		results = append(results, ecoData)
	}

	if err = rows.Err(); err != nil {
		log.Printf("RecommendEcosByUserEmbedding: row iteration error for userID '%d' - %v", userID, err)
		return nil, err
	}

	return results, nil
}

func (r *PostgresEcoRepo) UpdateLuv(ecoID, userID int64) (luvved bool, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return false, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var exists bool
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 FROM eco_luv_events
			WHERE eco_id = $1 AND user_id = $2
		);
	`
	if err = tx.QueryRow(checkQuery, ecoID, userID).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check luv existence: %w", err)
	}

	if exists {
		
		deleteQuery := `
			DELETE FROM eco_luv_events
			WHERE eco_id = $1 AND user_id = $2;
		`
		if _, err = tx.Exec(deleteQuery, ecoID, userID); err != nil {
			return false, fmt.Errorf("failed to delete luv event: %w", err)
		}

		updateQuery := `
			UPDATE eco_data
			SET luv_count = GREATEST(luv_count - 1, 0)
			WHERE eco_id = $1;
		`
		if _, err = tx.Exec(updateQuery, ecoID); err != nil {
			return false, fmt.Errorf("failed to decrement luv count: %w", err)
		}

		log.Printf("UpdateLuv: User %d removed luv from eco %d", userID, ecoID)
		luvved = false

	} else {
		
		insertQuery := `
			INSERT INTO eco_luv_events (eco_id, user_id)
			VALUES ($1, $2);
		`
		if _, err = tx.Exec(insertQuery, ecoID, userID); err != nil {
			return false, fmt.Errorf("failed to insert luv event: %w", err)
		}

		updateQuery := `
			UPDATE eco_data
			SET luv_count = luv_count + 1
			WHERE eco_id = $1;
		`
		if _, err = tx.Exec(updateQuery, ecoID); err != nil {
			return false, fmt.Errorf("failed to increment luv count: %w", err)
		}

		log.Printf("UpdateLuv: User %d added luv to eco %d", userID, ecoID)
		luvved = true
	}

	if err = tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return luvved, nil
}


func (r *PostgresEcoRepo) LuvStatus(ecoID, userID int64) (bool, int64, error) {
	var (
		luvved     bool
		totalLuvs  int64
	)

	// 1️⃣ Check if this user has liked this eco
	luvCheckQuery := `
        SELECT EXISTS (
            SELECT 1 FROM eco_luv_events
            WHERE eco_id = $1 AND user_id = $2
        );
    `
	if err := r.db.QueryRow(luvCheckQuery, ecoID, userID).Scan(&luvved); err != nil {
		return false, 0, fmt.Errorf("failed to check luv status: %w", err)
	}

	// 2️⃣ Get total luvs for this eco
	totalQuery := `
        SELECT COUNT(*)
        FROM eco_luv_events
        WHERE eco_id = $1;
    `
	if err := r.db.QueryRow(totalQuery, ecoID).Scan(&totalLuvs); err != nil {
		return false, 0, fmt.Errorf("failed to get total luvs: %w", err)
	}

	return luvved, totalLuvs, nil
}

func (r *PostgresEcoRepo) GetEcoMetaData(ecoID int64) (models.EcoMetaData, error) {
	query := `
		SELECT 
			e.eco_id,
			e.uploader_id,
			e.eco_text,
			e.eco_url,
			e.images_count,
			e.created_at,
			e.view_count,
			e.tags,
			e.luv_count,
			e.comment_count,
			u.user_profile_name,
			u.user_handle
		FROM eco_data e
		JOIN user_data_table u ON u.user_id = e.uploader_id
		WHERE e.eco_id = $1;
	`

	var ecoData models.EcoMetaData
	row := r.db.QueryRow(query, ecoID)

	err := row.Scan(
		&ecoData.Eco_Id,
		&ecoData.Uploader_ID,
		&ecoData.Eco_Text,
		&ecoData.Eco_Url,
		&ecoData.Images_Count,
		&ecoData.Created_At,
		&ecoData.View_Count,
		pq.Array(&ecoData.Tags),
		&ecoData.Luv_Count,
		&ecoData.Comment_Count,
		&ecoData.Uploader_Name,
		&ecoData.Uploader_Handle,
	)

	if err != nil {
		log.Printf("GetEcoMetaData: Failed to retrieve eco metadata for eco_id %d - %v", ecoID, err)
	} else {
		log.Printf("GetEcoMetaData: Successfully retrieved eco metadata for eco_id %d", ecoID)
	}

	return ecoData, err
}

func (r *PostgresEcoRepo) GetAllEcoIDs(limit, offset int) ([]int64, error) {
	query := `
		SELECT eco_id
		FROM eco_data
		ORDER BY eco_id
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		log.Println("GetAllEcoIDs: query failed -", err)
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Println("GetAllEcoIDs: scan failed -", err)
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

func (r *PostgresEcoRepo) GetEcoEngagementInfo(ecoID int64) (
	views, luvs, comments, shares int64, lastScore float64,
	timestampStr string,
	err error,
) {
	query := `
		SELECT 
			view_count,
			luv_count,
			comment_count,
			share_count,
			last_trending_score,
			created_at
		FROM eco_data
		WHERE eco_id = $1;
	`

	err = r.db.QueryRow(query, ecoID).Scan(
		&views,
		&luvs,
		&comments,
		&shares,
		&lastScore,
		&timestampStr,
	)
	if err != nil {
		log.Println("GetEcoEngagementInfo error for", ecoID, ":", err)
		return 0, 0, 0, 0, 0, "", err
	}

	return views, luvs, comments, shares, lastScore, timestampStr, nil
}

func (r *PostgresEcoRepo) SaveEcoTrendingDelta(ecoID int64, score float64, delta float64) error {
	query := `
		UPDATE eco_data
		SET 
			last_trending_score = $1,
			trending_delta = $2,
		WHERE eco_id = $3;
	`

	_, err := r.db.Exec(query, score, delta, ecoID)
	if err != nil {
		log.Println("SaveEcoTrendingDelta error for", ecoID, ":", err)
	}
	return err
}

func (r *PostgresEcoRepo) GetTrendingEcos(limit, offset int) ([]models.EcoMetaData, error) {
	query := `
		SELECT 
			e.eco_id,
			e.eco_text,
			e.eco_url,
			e.images_count,
			e.created_at,
			e.view_count,
			e.tags,
			e.luv_count,
			e.comment_count,
			u.user_profile_name AS uploader_name,
			u.user_handle AS uploader_handle,
			u.user_id
		FROM eco_data e
		JOIN user_data_table u ON e.uploader_id = u.user_id
		ORDER BY e.trending_delta DESC
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		log.Println("GetTrendingEcos query error:", err)
		return nil, err
	}
	defer rows.Close()

	var results []models.EcoMetaData
	for rows.Next() {
		var eco models.EcoMetaData
		err := rows.Scan(
			&eco.Eco_Id,
			&eco.Eco_Text,
			&eco.Eco_Url,
			&eco.Images_Count,
			&eco.Created_At,
			&eco.View_Count,
			pq.Array(&eco.Tags),
			&eco.Luv_Count,
			&eco.Comment_Count,
			&eco.Uploader_Name,
			&eco.Uploader_Handle,
			&eco.Uploader_ID,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, eco)
	}

	return results, rows.Err()
}

func (r *PostgresEcoRepo) GetEchoScore(echoID int64) (models.EchoScore, error) {
    echo_score := models.EchoScore{}
    query := `
        SELECT 
            AVG(quality),
            AVG(ai_usage),
            COUNT(quality),
            COUNT(ai_usage)
        FROM eco_votes
        WHERE eco_id = $1;
    `
    err := r.db.QueryRow(query, echoID).Scan(&echo_score.Echo_Quality, &echo_score.Echo_AI_Usage, &echo_score.Total_Qualtiy_Votes, &echo_score.Total_Ai_Votes)
    if err != nil{
        return echo_score, err
    }
    return echo_score, nil
}