/////  Video Meta data -> VMD

package repository

import (
	"GoServer/models"
	"database/sql"
    "fmt"
    "log"
)

// Interface
type VMDrepo interface {
    CreateNewVMD(vmd *models.VideoMetaData) error
    SearchVMDs(keyword string) ([]models.VideoSearchResult, error)
    GetSpecificVideoMD(video_id int64) ( models.VideoMetaData, error)
    VideoViewUpdate(video_id int64) error
    SearchVMDsBy(userID int64) ([]models.VideoMetaData, error)
}


// Postgres implementation
type PostgresVMDRepo struct {
    db *sql.DB
}

func NewPostgresVMDRepo(db *sql.DB) VMDrepo {
    return &PostgresVMDRepo{db: db}
}


/// CreateNewVMD for Postgres
func (r *PostgresVMDRepo) CreateNewVMD(vmd *models.VideoMetaData) error {
    query := `INSERT INTO video_data(uploader_id, title, video_info, video_url) VALUES($1, $2, $3, $4)`

    _, err := r.db.Exec(query,
        vmd.Uploader_ID,
        vmd.Title,
        vmd.Video_Info,
        vmd.Video_Url,
    )
    if err != nil {
        log.Printf("CreateNewVMD: Failed to insert video metadata for user %d, video %s - %v", vmd.Uploader_ID, vmd.Video_Url, err)
    } else {
        log.Printf("CreateNewVMD: Successfully created video metadata for user %d, video %s (title: %s)", vmd.Uploader_ID, vmd.Video_Url, vmd.Title)
    }
    return err
}



func (r *PostgresVMDRepo) SearchVMDs(keyword string) ([]models.VideoSearchResult, error) {
    query := `
    SELECT v.video_id, u.user_profile_name AS uploader_name, v.title, v.views, v.video_url, v.upload_time
    FROM video_data v
    JOIN user_data_table u ON v.uploader_id = u.user_id
    WHERE v.title ILIKE $1 OR v.video_info ILIKE $1
    `
    rows, err := r.db.Query(query, "%"+keyword+"%")
    if err != nil {
        log.Printf("SearchVMDs: Failed to execute search query for keyword '%s' - %v", keyword, err)
        return nil, err
    }
    defer rows.Close()

    var results []models.VideoSearchResult
    for rows.Next() {
        var res models.VideoSearchResult
        if err := rows.Scan(&res.VideoID, &res.UploaderName, &res.Title, &res.Views, &res.VideoURL, &res.Date); err != nil {
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


func (r *PostgresVMDRepo) GetSpecificVideoMD(video_id int64) (models.VideoMetaData, error){
    query := `
        SELECT 
            v.video_id,
            v.title,
            v.agrees,
            v.disagrees,
            v.video_info,
            v.upload_time,
            u.user_profile_name,
            u.user_handle
        FROM video_data v
        JOIN user_data_table u ON v.uploader_id = u.user_id
        WHERE v.video_id = $1;
        `
        var videoData models.VideoMetaData
        row := r.db.QueryRow(query, video_id)

        err := row.Scan(
            &videoData.Video_ID,
            &videoData.Title,
            &videoData.Agrees,
            &videoData.Disagrees,
            &videoData.Video_Info,
            &videoData.Upload_Time,
            &videoData.Uploader_Name,
            &videoData.Uploader_Handle,
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
                v.title,
                v.views,
                v.video_info,
                v.upload_time,
                v.video_url,
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
            &videoData.Title,
            &videoData.Views,
            &videoData.Video_Info,
            &videoData.Upload_Time,
            &videoData.Video_Url,
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