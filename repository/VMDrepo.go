/////  Video Meta data -> VMD

package repository

import (
	"GoServer/models"
	"database/sql"
    "fmt"
    "log"
    "github.com/lib/pq"
)

// Interface
type VMDrepo interface {
    CreateNewVMD(vmd *models.VideoMetaData) (int64, error)
    SearchVMDs(keyword string, limit, offset int) ([]models.VideoSearchResult, error)
    GetSpecificVideoMD(video_id int64, user_id int64) ( models.VideoMetaData, error)
    VideoViewUpdate(video_id int64) error
    SearchVMDsBy(userID int64) ([]models.VideoMetaData, error)
    SimilaritySearch(video_id int64, limit int , offset int) ([]models.VideoMetaData, error)
    UpdateLuv(videoID int64, user_id int64) (bool ,error)
    GetUserWatchHistory(user_id int64,limit int, offset int)([]models.VideoMetaData, error)
    UpdateUserEmbeddingsFromVideoHistory(userID int64, historyCount int) error
    RecommendVideosByUserEmbedding(userID int64, limit, offset int) ([]models.VideoMetaData, error)
    GetSavedVideos(userID int64, limit, offset int) ([]models.VideoMetaData, error)
    GetEngagementInfo(videoID int64)  (int64, int64, int64, int64, float64, string, error)
    SaveTrendingDelta(videoID int64, currentTrendingScore float64, trendingDelta float64) error
    GetVideoIDsPaginated(limit, offset int) ([]int64, error)
    GetTrendingVMDsPaginated(userID int64, limit, offset int) ([]models.VideoMetaData, error)
    DeleteMyHistory(userID int64) error
    GetVideoScore(videoID int64) (models.VideoScore, error)
}


// Postgres implementation
type PostgresVMDRepo struct {
    db *sql.DB
}

func NewPostgresVMDRepo(db *sql.DB) VMDrepo {
    return &PostgresVMDRepo{db: db}
}


/// CreateNewVMD for Postgres
func (r *PostgresVMDRepo) CreateNewVMD(vmd *models.VideoMetaData) (int64, error) {
    query := `INSERT INTO video_data(uploader_id, title, video_info, video_url, tags) VALUES($1, $2, $3, $4, $5) RETURNING video_id`

    var videoID int64
    err := r.db.QueryRow(query,
        vmd.Uploader_ID,
        vmd.Title,
        vmd.Video_Info,
        vmd.Video_Url,
        pq.Array(vmd.Tags),
    ).Scan(&videoID)
    
    if err != nil {
        log.Printf("CreateNewVMD: Failed to insert video metadata for user %d, video %s - %v", vmd.Uploader_ID, vmd.Video_Url, err)
        return 0, err
    }
    
    log.Printf("CreateNewVMD: Successfully created video metadata for user %d, video %s (title: %s, video_id: %d)", vmd.Uploader_ID, vmd.Video_Url, vmd.Title, videoID)
    return videoID, nil
}

func (r *PostgresVMDRepo) SearchVMDs(keyword string, limit, offset int) ([]models.VideoSearchResult, error) {
    query := `
    SELECT v.video_id, u.user_profile_name AS uploader_name, u.user_handle AS uploader_handle, v.title, v.views, v.video_url, v.upload_time, v.tags
    FROM video_data v
    JOIN user_data_table u ON v.uploader_id = u.user_id
    WHERE v.title ILIKE $1 OR v.video_info ILIKE $1 OR array_to_string(v.tags, ' ') ILIKE $1
    LIMIT $2 OFFSET $3
    `
    rows, err := r.db.Query(query, "%"+ keyword +"%", limit, offset)
    if err != nil {
        log.Printf("SearchVMDs: Failed to execute search query for keyword '%s' - %v", keyword, err)
        return nil, err
    }
    defer rows.Close()

    var results []models.VideoSearchResult
    for rows.Next() {
        var res models.VideoSearchResult
        if err := rows.Scan(&res.VideoID, &res.UploaderName, &res.UploaderHandle, &res.Title, &res.Views, &res.VideoURL, &res.Date, pq.Array(&res.Tags)); err != nil {
            log.Printf("SearchVMDs: Failed to scan row for keyword '%s' - %v", keyword, err)
            return nil, err
        }
        results = append(results, res)
    }

    if err := rows.Err(); err != nil {
        log.Printf("SearchVMDs: Row iteration error for keyword '%s' - %v", keyword, err)
        return nil, err
    }

    log.Printf("SearchVMDs: Found %d results for keyword '%s'", len(results), keyword)
    return results, nil

}
func (r *PostgresVMDRepo) VideoViewUpdate(videoID int64) error {
    query := `UPDATE video_data SET views = views + 1 WHERE video_id = $1`

    _, err := r.db.Exec(query, videoID)
    if err != nil {
        log.Printf("VideoViewUpdate: Failed to update views for video_id %d - %v", videoID, err)
        return fmt.Errorf("failed to update views for video_id %d: %w", videoID, err)
    }

    log.Printf("VideoViewUpdate: Successfully incremented views for video_id %d", videoID)
    return nil
}


func (r *PostgresVMDRepo) GetSpecificVideoMD(video_id int64, user_id int64) (models.VideoMetaData, error){
    query := `
        SELECT 
            v.video_id,
            v.uploader_id,
            v.title,
            v.luv,
            v.video_info,
            v.upload_time,
            u.user_profile_name,
            u.user_handle,
            v.views,
            v.tags,
            CASE WHEN vle.user_id IS NOT NULL THEN true ELSE false END AS already_luved
        FROM video_data v
        JOIN user_data_table u ON v.uploader_id = u.user_id
        LEFT JOIN video_luv_events vle ON v.video_id = vle.video_id AND vle.user_id = $2
        WHERE v.video_id = $1;
        `
        var videoData models.VideoMetaData
        row := r.db.QueryRow(query, video_id, user_id)

        err := row.Scan(
            &videoData.Video_ID,
            &videoData.Uploader_ID,
            &videoData.Title,
            &videoData.Luvs,
            &videoData.Video_Info,
            &videoData.Upload_Time,
            &videoData.Uploader_Name,
            &videoData.Uploader_Handle,
            &videoData.Views,
            pq.Array(&videoData.Tags),
            &videoData.Already_Luved,
        )   

        if err != nil {
            log.Printf("GetSpecificVideoMD: Failed to retrieve video metadata for video_id %d - %v", video_id, err)
        } else {
            log.Printf("GetSpecificVideoMD: Successfully retrieved video metadata for video_id %d (title: %s)", video_id, videoData.Title)
        }

        return videoData , err 
}


func (r *PostgresVMDRepo) SearchVMDsBy(userID int64) ([]models.VideoMetaData, error) {
    query := `
                SELECT 
                v.video_id,
                u.user_id,
                v.title,
                v.views,
                v.video_info,
                v.upload_time,
                v.video_url,
                v.tags,
                u.user_profile_name AS uploader_name,
                u.user_handle AS uploader_handle
                FROM video_data v
                JOIN user_data_table u
                ON v.uploader_id = u.user_id
                WHERE v.uploader_id = $1;
    `
    rows, err := r.db.Query(query, userID)
    if err != nil {
        log.Printf("SearchVMDsBy: Failed to execute search query for userID '%d' - %v", userID, err)
        return nil, err
    }
    defer rows.Close()  
    var results []models.VideoMetaData;
    for rows.Next(){
        var videoData models.VideoMetaData
        if err := rows.Scan(
            &videoData.Video_ID,
            &videoData.Uploader_ID,
            &videoData.Title,
            &videoData.Views,
            &videoData.Video_Info,
            &videoData.Upload_Time,
            &videoData.Video_Url,
            pq.Array(&videoData.Tags),
            &videoData.Uploader_Name,
            &videoData.Uploader_Handle,
        ); err != nil {
            log.Printf("SearchVMDsBy: Failed to scan row for userID '%d' - %v", userID, err)
            return nil, err
        }
        results = append(results, videoData)
    }
    return results, nil
}


func (r *PostgresVMDRepo) SimilaritySearch(video_id int64, limit int , offset int) ([]models.VideoMetaData, error){

    query := `
        SELECT 
        v.video_id,
        u.user_id,
        u.user_profile_name AS uploader_name,
        u.user_handle AS uploader_handle,
        v.title,
        v.views,
        v.video_url,
        v.tags,
        v.upload_time
        FROM video_data v
        JOIN user_data_table u ON v.uploader_id = u.user_id
        WHERE v.video_id <> $1
        ORDER BY v.embeddings <=> (
        SELECT embeddings FROM video_data WHERE video_id = $1
        )
        LIMIT $2 OFFSET $3;
    `

	rows, err := r.db.Query(query, video_id, limit, offset)
	if err != nil {
		log.Printf("SimilaritySearch: failed for video_id %d - %v", video_id, err)
		return nil, err
	}
	defer rows.Close()

	var results []models.VideoMetaData
	for rows.Next() {
		var v models.VideoMetaData
		err := rows.Scan(&v.Video_ID,&v.Uploader_ID, &v.Uploader_Name, &v.Uploader_Handle, &v.Title, &v.Views, &v.Video_Url, pq.Array(&v.Tags), &v.Upload_Time)
		if err != nil {
			return nil, err
		}
		results = append(results, v)
	}

	return results, nil
}

func (r *PostgresVMDRepo) UpdateLuv(videoID, userID int64) (luvved bool, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return false, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Step 1: Check if the user already liked the video
	var exists bool
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 FROM video_luv_events
			WHERE video_id = $1 AND user_id = $2
		);
	`
	if err = tx.QueryRow(checkQuery, videoID, userID).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check luv existence: %w", err)
	}

	if exists {
		// Step 2a: User already liked → remove luv and decrement count
		deleteQuery := `
			DELETE FROM video_luv_events
			WHERE video_id = $1 AND user_id = $2;
		`
		if _, err = tx.Exec(deleteQuery, videoID, userID); err != nil {
			return false, fmt.Errorf("failed to delete luv event: %w", err)
		}

		updateQuery := `
			UPDATE video_data
			SET luv = GREATEST(luv - 1, 0)
			WHERE video_id = $1;
		`
		if _, err = tx.Exec(updateQuery, videoID); err != nil {
			return false, fmt.Errorf("failed to decrement luv count: %w", err)
		}

		log.Printf("UpdateLuv: User %d removed luv from video %d", userID, videoID)
		luvved = false

	} else {
		// Step 2b: User has not liked → add luv and increment count
		insertQuery := `
			INSERT INTO video_luv_events (video_id, user_id)
			VALUES ($1, $2);
		`
		if _, err = tx.Exec(insertQuery, videoID, userID); err != nil {
			return false, fmt.Errorf("failed to insert luv event: %w", err)
		}

		updateQuery := `
			UPDATE video_data
			SET luv = luv + 1
			WHERE video_id = $1;
		`
		if _, err = tx.Exec(updateQuery, videoID); err != nil {
			return false, fmt.Errorf("failed to increment luv count: %w", err)
		}

		log.Printf("UpdateLuv: User %d added luv to video %d", userID, videoID)
		luvved = true
	}

	// Step 3: Commit the transaction
	if err = tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return luvved, nil
}

func (r *PostgresVMDRepo) GetUserWatchHistory(user_id int64, limit, offset int) ([]models.VideoMetaData, error) {
	query := `
        SELECT 
            v.video_id,
            u.user_profile_name AS uploader_name,
            v.title,
            v.views,
            v.video_url,
            v.tags,
            v.upload_time,
            u.url,
            u.user_id,
            u.user_handle,
            h.watched_at
        FROM video_history_table h
        JOIN video_data v ON h.video_id = v.video_id
        JOIN user_data_table u ON v.uploader_id = u.user_id
        WHERE h.watcher_id = $1
        ORDER BY h.watched_at DESC
        LIMIT $2 OFFSET $3;
    `

	rows, err := r.db.Query(query, user_id, limit, offset)
	if err != nil {
		log.Printf("GetUserWatchHistory: query failed for user_id %d - %v", user_id, err)
		return nil, err
	}
	defer rows.Close()

	var history []models.VideoMetaData
	for rows.Next() {
		var v models.VideoMetaData
		err := rows.Scan(&v.Video_ID, &v.Uploader_Name, &v.Title, &v.Views, &v.Video_Url, pq.Array(&v.Tags), &v.Upload_Time, &v.Uploader_Url, &v.Uploader_ID, &v.Uploader_Handle, &v.Watched_At)
		if err != nil {
			return nil, err
		}
		history = append(history, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}


func (r *PostgresVMDRepo) UpdateUserEmbeddingsFromVideoHistory(userID int64, historyCount int) error {
    query := `
        UPDATE user_data_table
        SET embeddings = sub.mean_embedding
        FROM (
            SELECT watcher_id, avg(recent.embeddings) AS mean_embedding
            FROM (
                SELECT h.watcher_id, v.embeddings
                FROM video_history_table h
                JOIN video_data v ON h.video_id = v.video_id
                WHERE h.watcher_id = $1
                ORDER BY h.watched_at DESC
                LIMIT $2
            ) AS recent
            GROUP BY watcher_id
        ) AS sub
        WHERE user_data_table.user_id = sub.watcher_id;
    `

    _, err := r.db.Exec(query, userID, historyCount)
    if err != nil {
        log.Printf("Failed to update user embedding for user_id %d: %v", userID, err)
        return err
    }

    log.Printf("User embedding updated for user_id %d (last %d videos)", userID, historyCount)
    return nil
}

func (repo *PostgresVMDRepo) RecommendVideosByUserEmbedding(userID int64, limit, offset int) ([]models.VideoMetaData, error) {
	query := `
		SELECT 
			v.video_id,
			v.uploader_id,
			u.user_profile_name AS uploader_name,
			u.user_handle AS uploader_handle,
			v.title,
			v.views,
			v.video_url,
			v.tags,
			v.upload_time
		FROM video_data v
		JOIN user_data_table u ON v.uploader_id = u.user_id
		JOIN user_data_table usr ON usr.user_id = $1
		WHERE v.video_id NOT IN (
			SELECT video_id FROM video_history_table WHERE watcher_id = $1
		)
		ORDER BY usr.embeddings <=> v.embeddings ASC
		LIMIT $2 OFFSET $3;
	`

	rows, err := repo.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []models.VideoMetaData
	for rows.Next() {
		var v models.VideoMetaData
		err := rows.Scan(&v.Video_ID, &v.Uploader_ID, &v.Uploader_Name, &v.Uploader_Handle, &v.Title, &v.Views, &v.Video_Url, pq.Array(&v.Tags), &v.Upload_Time)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}

	return videos, nil
}

func (r *PostgresVMDRepo) GetSavedVideos(userID int64, limit, offset int) ([]models.VideoMetaData, error) {
	query := `
		SELECT 
			v.video_id,
			v.uploader_id,
			u.user_profile_name AS uploader_name,
			u.user_handle AS uploader_handle,
			v.title,
			v.views,
			v.video_url,
			v.tags,
			v.upload_time
		FROM saved_videos sv
		JOIN video_data v ON sv.video_id = v.video_id
		JOIN user_data_table u ON v.uploader_id = u.user_id
		WHERE sv.user_id = $1
		ORDER BY sv.created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []models.VideoMetaData
	for rows.Next() {
		var v models.VideoMetaData
		err := rows.Scan(&v.Video_ID, &v.Uploader_ID, &v.Uploader_Name, &v.Uploader_Handle, &v.Title, &v.Views, &v.Video_Url, pq.Array(&v.Tags), &v.Upload_Time)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}

	return videos, nil
}


func (r *PostgresVMDRepo) GetEngagementInfo(videoID int64) (int64, int64, int64, int64, float64, string, error){
    selectQuery := `
        SELECT views, luv, comments, shares, last_trending_score, upload_time
        FROM video_data
        WHERE video_id = $1;
    `

    var views, luv, comments, shares int64
    var last_trending_score float64
    var timeStamp string
    row := r.db.QueryRow(selectQuery, videoID)
    err := row.Scan(&views, &luv, &comments, &shares, &last_trending_score, &timeStamp)
    if err != nil {
        log.Printf("GetEngagementInfo: Failed to retrieve engagement info for video_id %d - %v", videoID, err)
        return 0, 0, 0, 0, 0, "", err
    }

    log.Printf("GetEngagementInfo: Successfully retrieved engagement info for video_id %d", videoID)
    return views, luv, comments, shares, last_trending_score, timeStamp, nil
}

func (repo *PostgresVMDRepo) SaveTrendingDelta(videoID int64, currentTrendingScore float64, trendingDelta float64) error {

    query := `
        UPDATE video_data
        SET trending_delta = $1,
            last_trending_score = $2,
            last_trending_updated_at = NOW()
        WHERE video_id = $3;
    `

    _, err := repo.db.Exec(query, trendingDelta, currentTrendingScore, videoID)
    return err
}

func (repo *PostgresVMDRepo) GetVideoIDsPaginated(limit, offset int) ([]int64, error) {
    rows, err := repo.db.Query(`
        SELECT video_id 
        FROM video_data
        ORDER BY video_id
        LIMIT $1 OFFSET $2
    `, limit, offset)

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

func (r *PostgresVMDRepo) GetTrendingVMDsPaginated(userID int64, limit, offset int) ([]models.VideoMetaData, error) {
    
    query := `
        SELECT 
            v.video_id,
            v.uploader_id,
            u.user_profile_name AS uploader_name,
            u.user_handle AS uploader_handle,
            v.title,
            v.views,
            v.video_url,
            v.tags,
            v.upload_time,
            u.user_handle
        FROM video_data v
        JOIN user_data_table u 
            ON v.uploader_id = u.user_id
        JOIN user_data_table usr 
            ON usr.user_id = $1
        WHERE v.video_id NOT IN (
            SELECT video_id 
            FROM video_history_table 
            WHERE watcher_id = $1
        )
        ORDER BY v.trending_delta DESC
        LIMIT $2 OFFSET $3;
    `

    rows, err := r.db.Query(query,userID,  limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    vmds := []models.VideoMetaData{}

    for rows.Next() {
        var v models.VideoMetaData

        err := rows.Scan(&v.Video_ID, &v.Uploader_ID, &v.Uploader_Name, &v.Uploader_Handle, &v.Title, &v.Views, &v.Video_Url, pq.Array(&v.Tags), &v.Upload_Time, &v.Uploader_Handle)
        if err != nil {
            return nil, err
        }

        vmds = append(vmds, v)
    }

    return vmds, nil
}

func (r *PostgresVMDRepo) DeleteMyHistory(userID int64) error {
    query := `
        DELETE FROM video_history_table
        WHERE watcher_id = $1;
    `
    _, err := r.db.Exec(query, userID)
    return err
}

func (r *PostgresVMDRepo) GetVideoScore(videoID int64) (models.VideoScore, error) {
    video_score := models.VideoScore{}
    query := `
        SELECT 
            AVG(quality),
            AVG(ai_usage),
            COUNT(quality),
            COUNT(ai_usage)
        FROM video_votes
        WHERE video_id = $1;
    `
    err := r.db.QueryRow(query, videoID).Scan(&video_score.Video_Quality, &video_score.Video_AI_Usage, &video_score.Total_Qualtiy_Votes, &video_score.Total_Ai_Votes)
    if err != nil{
        return video_score, err
    }
    return video_score, nil
}