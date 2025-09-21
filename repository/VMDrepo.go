/////  Video Meta data -> VMD

package repository

import (
    "GoServer/models"
    "database/sql"
)

// Interface
type VMDrepo interface {
    CreateNewVMD(vmd *models.VideoMetaData) error
    SearchVMDs(keyword string) ([]models.VideoSearchResult, error)
    GetSpecificVideoMD(video_id int64) ( models.VideoMetaData, error)
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
        return nil, err
    }
    defer rows.Close()

    var results []models.VideoSearchResult
    for rows.Next() {
        var res models.VideoSearchResult
        if err := rows.Scan(&res.VideoID, &res.UploaderName, &res.Title, &res.Views, &res.VideoURL, &res.Date); err != nil {
            return nil, err
        }
        results = append(results, res)
    }

    return results, rows.Err()

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



        return videoData , err 
}