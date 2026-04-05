package repository

import (
	"GoServer/models"
	"database/sql"
	"fmt"
	"log"
	"strings"

)

type UserRepo interface {
    CreateNewUser(user *models.UserData, auth *models.UserAuth) error
	GetUser(userID int64) (*models.UserData, error)
	UpadateUserProfile(user *models.UserData) error
    AddUserAuth(auth *models.UserAuth) error
	CheckUser(userHandle string) (int64, string, error)
	CheckEmailExists(email string) (bool, error)
	FollowUser(FollowerID int64, FolloweeID int64) (error)
	UnfollowUser(FollowerID int64, FolloweeID int64) (error)
	GetAllFollowers(FolloweeID int64) ([]int64,error)
	GetAllFollowees(followerID int64) ([]int64 ,error)
	AddConnection(requesterID int64, respondentID int64) (error)
	SearchWithKeyword(keyword string) ([]int64, error)
	AddVideoInUserHistory(userID int64, videoID int64)(error)
	AllUsersReturn(limit int, offset int) ([]int64, error)
	GetFollowingInfo(userID, requesterID int64)(models.FollowData, error)
	UserSavedEco(userID int64, ecoID int64)(bool, error)
	UserSavedVideo(userID int64, videoID int64) (bool, error)
	UserEcoSavedStatus(userID int64, ecoID int64)(bool, error)
	UserVideoSavedStatus(userID int64, videoID int64) (bool, error)
	GetTurbomaxStatusOfUser(userID int64)(bool, error)
	GetUserAnalytics(userID int64)(map[string]map[string]int, error)
	UpsertVideoVote(videoID, userID int64, quality, aiUsage int) error
	UpsertEcoVote(ecoID, userID int64, quality, aiUsage int) error
}


// Postgres implementation
type PostgresUserRepo struct {
    db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) UserRepo {
    return &PostgresUserRepo{db: db}
}


func (r *PostgresUserRepo) CreateNewUser(user *models.UserData, auth *models.UserAuth) error {
    
    // Starting a new Transaction.
    tx, err := r.db.Begin()
	if err != nil {
		log.Println("Error starting the db transaction")
        return err
	}


    // If the func returns an error, rollback the transaction
    defer func() {
        if err != nil {
            _ = tx.Rollback()
        }
    }()


    userQuery := `
		INSERT INTO user_data_table (user_handle, user_profile_name, user_description, from_location,  gender, url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING user_id
	`
	err = tx.QueryRow(userQuery,
		user.UserHandle,
		user.UserProfileName,
		user.UserDescription,
		user.FromLocation,
		user.Gender,
		user.Url,
	).Scan(&user.UserID)
	if err != nil {
        log.Printf("Error inserting into user_data_table: %v", err)
         return fmt.Errorf("error inserting into user_data_table table", err)
	}



    authQuery := `
		INSERT INTO user_authentication (user_id, user_login_account, user_phone_number, user_hashed_password, account_created_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		RETURNING auth_id
	`
	err = tx.QueryRow(authQuery,
		user.UserID,
		auth.UserLoginAccount,
		auth.UserPhoneNumber,
		auth.UserHashedPassword,
	).Scan(&auth.AuthID)
	if err != nil {
		
        return fmt.Errorf("error inserting into user_authentication table", err)
	}


    auth.UserID = user.UserID

	// Comminting the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("failed to commit transaction")
        return err
	}

	return nil
}






func (r *PostgresUserRepo) UpadateUserProfile(user *models.UserData) error {

	// some authentication Check sould be here
	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if user.UserHandle != "" {
		setClauses = append(setClauses, fmt.Sprintf("user_handle = $%d", argPos))
		args = append(args, user.UserHandle)
		argPos++
	}
	if user.UserProfileName != "" {
		setClauses = append(setClauses, fmt.Sprintf("user_profile_name = $%d", argPos))
		args = append(args, user.UserProfileName)
		argPos++
	}
	if user.UserDescription != "" {
		setClauses = append(setClauses, fmt.Sprintf("user_description = $%d", argPos))
		args = append(args, user.UserDescription)
		argPos++
	}
	if user.FromLocation != "" {
		setClauses = append(setClauses, fmt.Sprintf("from_location = $%d", argPos))
		args = append(args, user.FromLocation)
		argPos++
	}
	if user.UserDateOfBirth != "" {
		setClauses = append(setClauses, fmt.Sprintf("user_date_of_birth = $%d", argPos))
		args = append(args, user.UserDateOfBirth)
		argPos++
	}
	if user.Gender != "" {
		setClauses = append(setClauses, fmt.Sprintf("gender = $%d", argPos))
		args = append(args, user.Gender)
		argPos++
	}

	// No fields provided? -> do nothing
	if len(setClauses) == 0 {
		return nil
	}

	query := fmt.Sprintf(
		"UPDATE user_data_table SET %s WHERE user_id = $%d",
		strings.Join(setClauses, ", "),
		argPos,
	)
	args = append(args, user.UserID)

	_, err := r.db.Exec(query, args...)
	return err
}





/// adding Another authentication method for same uesr
func (r *PostgresUserRepo) AddUserAuth(auth *models.UserAuth) error {
	query := `
		INSERT INTO user_authentication (user_id, user_login_account, user_phone_number, user_hashed_password, account_created_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		RETURNING auth_id
	`
	return r.db.QueryRow(query,
		auth.UserID,
		auth.UserLoginAccount,
		auth.UserPhoneNumber,
		auth.UserHashedPassword,
	).Scan(&auth.AuthID)
}


func (r *PostgresUserRepo) CheckUser(userHandle string) (int64, string, error){
	var userID int64
	var passwordHash string
	err := r.db.QueryRow(
		`
		SELECT 
			d.user_id,
			a.user_hashed_password
		FROM 
			user_data_table d
		JOIN 
			user_authentication a 
		ON 
			d.user_id = a.user_id
		WHERE 
			d.user_handle = $1;
		`, userHandle).Scan(&userID, &passwordHash)
	if err != nil {
		return 0 , "" ,fmt.Errorf("failed to get user handle: %v", err)
	}

	return userID, passwordHash, nil
}

func (r *PostgresUserRepo) CheckEmailExists(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM user_authentication WHERE user_login_account = $1)`,
		email,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return exists, nil
}



func (r *PostgresUserRepo)FollowUser(FollowerID int64, FolloweeID int64) (error) {
	query := `
	INSERT INTO follow_table (follower_id, followee_id)
	VALUES ($1, $2)
	ON CONFLICT (follower_id, followee_id) DO NOTHING;
	`
	_, err := r.db.Exec(query, FollowerID, FolloweeID)
	return err
}

func (r *PostgresUserRepo)UnfollowUser(FollowerID int64, FolloweeID int64) (error) {
	query := `
		DELETE FROM follow_table
		WHERE follower_id = $1 AND followee_id = $2;
	`
	_, err := r.db.Exec(query, FollowerID, FolloweeID)
	return err
}


func (r *PostgresUserRepo)GetUser(UserID int64) (*models.UserData, error) {
	query := `
		SELECT 
			user_handle, 
			user_description,
			user_profile_name, 
			from_location,
			gender

		FROM user_data_table
		WHERE user_id = $1
	`

	row := r.db.QueryRow(query, UserID)

	var user models.UserData
	err := row.Scan(
		&user.UserHandle,
		&user.UserDescription,
		&user.UserProfileName,
		&user.FromLocation,
		&user.Gender,
	)
	user.UserID = UserID
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepo)GetAllFollowers(followeeID int64) ([]int64 ,error) {
	query := `
        SELECT follower_id
        FROM follow_table
        WHERE followee_id = $1;
    `
	rows, err := r.db.Query(query, followeeID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var followers []int64
    for rows.Next() {
        var followerID int64
        if err := rows.Scan(&followerID); err != nil {
            return nil, err
        }
        followers = append(followers, followerID)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return followers, nil
}

func (r *PostgresUserRepo)GetAllFollowees(followerID int64) ([]int64 ,error) {
	query := `
        SELECT followee_id
        FROM follow_table
        WHERE follower_id = $1;
    `
	rows, err := r.db.Query(query, followerID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var followees []int64
    for rows.Next() {
        var followeeID int64
        if err := rows.Scan(&followeeID); err != nil {
            return nil, err
        }
        followees = append(followees, followeeID)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return followees, nil
}

func (r *PostgresUserRepo)SearchWithKeyword(keyword string) ([]int64, error) {
	query := `  SELECT user_id
				FROM user_data_table
				WHERE 
				user_handle ILIKE '%' || $1 || '%' OR
				user_description ILIKE '%' || $1 || '%' OR
				from_location ILIKE '%' || $1 || '%' OR
				user_profile_name ILIKE '%' || $1 || '%'
				ORDER BY 
				CASE
					WHEN user_handle ILIKE '%' || $1 || '%' THEN 1
					WHEN user_profile_name ILIKE '%' || $1 || '%' THEN 2
					ELSE 3
				END ASC;
		`
	rows, err := r.db.Query(query, keyword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDList []int64
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDList = append(userIDList, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIDList, nil
}

func (r *PostgresUserRepo)AddConnection(user1ID int64, user2ID int64) (error){

	query := `
    INSERT INTO connections (user1_id, user2_id)
    VALUES (LEAST($1::BIGINT, $2::BIGINT), GREATEST($1::BIGINT, $2::BIGINT))
    ON CONFLICT (user1_id, user2_id) DO NOTHING;
`
	_, err := r.db.Exec(query, user1ID, user2ID)
	return err
}

func (r *PostgresUserRepo)AddVideoInUserHistory(userID int64, videoID int64)(error){
	query := `
		INSERT INTO video_history_table (watcher_id, video_id)
		VALUES ($1, $2)
	`
	_, err := r.db.Exec(query, userID, videoID)
	if err != nil {
		return fmt.Errorf("failed to insert video history: %w", err)
	}
	return nil
}



func (r *PostgresUserRepo) AllUsersReturn(limit int, offset int) ([]int64, error) {
	query := `
        SELECT user_id
        FROM user_data_table
        ORDER BY user_id ASC
        LIMIT $1 OFFSET $2;
    `

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		log.Printf("AllUsersReturn: query failed - %v", err)
		return nil, err
	}
	defer rows.Close()

	var userIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}

func (r *PostgresUserRepo) GetFollowingInfo(userID, requesterID int64) (models.FollowData, error) {
	var data models.FollowData

	followerQuery := `SELECT COUNT(*) FROM follow_table WHERE followee_id = $1`
	if err := r.db.QueryRow(followerQuery, userID).Scan(&data.FollowerCount); err != nil {
		return data, fmt.Errorf("failed to get follower count: %w", err)
	}

	followeeQuery := `SELECT COUNT(*) FROM follow_table WHERE follower_id = $1`
	if err := r.db.QueryRow(followeeQuery, userID).Scan(&data.FolloweeCount); err != nil {
		return data, fmt.Errorf("failed to get followee count: %w", err)
	}

	if requesterID > 0 && requesterID != userID {
		existsQuery := `
			SELECT EXISTS(
				SELECT 1 
				FROM follow_table 
				WHERE follower_id = $1 AND followee_id = $2
			)`
		if err := r.db.QueryRow(existsQuery, requesterID, userID).Scan(&data.AlreadyFollowed); err != nil {
			return data, fmt.Errorf("failed to check follow status: %w", err)
		}
	}

	return data, nil
}

func (r *PostgresUserRepo)UserSavedVideo(userID int64, videoID int64)(bool, error){
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM saved_videos WHERE user_id = $1 AND video_id = $2)`, userID, videoID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check video existence: %w", err)
	}
	
	if exists {
		_, err = r.db.Exec(`DELETE FROM saved_videos WHERE user_id = $1 AND video_id = $2`, userID, videoID)
		if err != nil {
			return false, fmt.Errorf("failed to delete video: %w", err)
		}
		_, err = r.db.Exec(`UPDATE video_data SET saves_count = GREATEST(saves_count - 1, 0) WHERE video_id = $1`, videoID)
		if err != nil {
			return false, fmt.Errorf("failed to decrement saves_count: %w", err)
		}
		return false, nil
	} else {
		_, err = r.db.Exec(`INSERT INTO saved_videos (user_id, video_id) VALUES ($1, $2)`, userID, videoID)
		if err != nil {
			return false, fmt.Errorf("failed to save video: %w", err)
		}
		_, err = r.db.Exec(`UPDATE video_data SET saves_count = saves_count + 1 WHERE video_id = $1`, videoID)
		if err != nil {
			return false, fmt.Errorf("failed to increment saves_count: %w", err)
		}
		return true, nil
	}
}

func (r *PostgresUserRepo)UserSavedEco(userID int64, ecoID int64)(bool, error){
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM saved_ecos WHERE user_id = $1 AND eco_id = $2)`, userID, ecoID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check eco existence: %w", err)
	}
	
	if exists {
		_, err = r.db.Exec(`DELETE FROM saved_ecos WHERE user_id = $1 AND eco_id = $2`, userID, ecoID)
		if err != nil {
			return false, fmt.Errorf("failed to delete eco: %w", err)
		}
		_, err = r.db.Exec(`UPDATE eco_data SET saves_count = GREATEST(saves_count - 1, 0) WHERE eco_id = $1`, ecoID)
		if err != nil {
			return false, fmt.Errorf("failed to decrement saves_count: %w", err)
		}
		return false, nil
	} else {
		_, err = r.db.Exec(`INSERT INTO saved_ecos (user_id, eco_id) VALUES ($1, $2)`, userID, ecoID)
		if err != nil {
			return false, fmt.Errorf("failed to save eco: %w", err)
		}
		_, err = r.db.Exec(`UPDATE eco_data SET saves_count = saves_count + 1 WHERE eco_id = $1`, ecoID)
		if err != nil {
			return false, fmt.Errorf("failed to increment saves_count: %w", err)
		}
		return true, nil
	}
}

func (r *PostgresUserRepo)UserVideoSavedStatus(userID int64, videoID int64)(bool, error){
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM saved_videos WHERE user_id = $1 AND video_id = $2)`, userID, videoID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check video existence: %w", err)
	}
	
	return exists, nil
}

func (r *PostgresUserRepo)UserEcoSavedStatus(userID int64, ecoID int64)(bool, error){
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM saved_ecos WHERE user_id = $1 AND eco_id = $2)`, userID, ecoID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check eco existence: %w", err)
	}
	
	return exists, nil
}

func (r * PostgresUserRepo)GetTurbomaxStatusOfUser(userID int64)(bool, error){
	var subscription_active bool
	query := `
		SELECT 
			(is_active = TRUE AND expiry_date > NOW()) AS active
		FROM turbomax_status_table
		WHERE user_id = $1
		ORDER BY turbomax_id DESC
		LIMIT 1;
`
	err := r.db.QueryRow(query, userID).Scan(&subscription_active)
	if err == sql.ErrNoRows{
		return false , nil
	}
	if err!= nil{
		return false , err
	}
	return subscription_active, nil

}


func (r *PostgresUserRepo) GetUserAnalytics(userID int64) (map[string]map[string]int, error) {
	result := make(map[string]map[string]int)
	result["VideoUploads"] = make(map[string]int)
	result["EcoUploads"] = make(map[string]int)

	queryVideoUploads := `
		SELECT
			TO_CHAR(upload_time::date, 'YYYY-MM-DD') AS day,
			COUNT(*) AS count
		FROM video_data
		WHERE uploader_id = $1
		GROUP BY upload_time::date
		ORDER BY upload_time::date;`

	rows, err := r.db.Query(queryVideoUploads, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var day string
		var count int
		if err := rows.Scan(&day, &count); err != nil {
			return nil, err
		}
		result["VideoUploads"][day] = count
	}

	queryEcoUploads := `
		SELECT
			TO_CHAR(created_at::date, 'YYYY-MM-DD') AS day,
			COUNT(*) AS count
		FROM eco_data
		WHERE uploader_id = $1
		GROUP BY created_at::date
		ORDER BY created_at::date;`

	rows, err = r.db.Query(queryEcoUploads, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var day string
		var count int
		if err := rows.Scan(&day, &count); err != nil {
			return nil, err
		}
		result["EcoUploads"][day] = count
	}

	return result, nil
}

func (r *PostgresUserRepo) UpsertVideoVote(
	videoID, userID int64,
	quality, aiUsage int,
) error {

	query := `
	INSERT INTO video_votes (video_id, user_id, quality, ai_usage)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (video_id, user_id)
	DO UPDATE SET
		quality = EXCLUDED.quality,
		ai_usage = EXCLUDED.ai_usage,
		updated_at = now();
	`

	_, err := r.db.Exec(query, videoID, userID, quality, aiUsage)
	return err
}

func (r *PostgresUserRepo) UpsertEcoVote(
	ecoID, userID int64,
	quality, aiUsage int,
) error {

	query := `
	INSERT INTO eco_votes (eco_id, user_id, quality, ai_usage)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (eco_id, user_id)
	DO UPDATE SET
		quality = EXCLUDED.quality,
		ai_usage = EXCLUDED.ai_usage,
		updated_at = now();
	`

	_, err := r.db.Exec(query, ecoID, userID, quality, aiUsage)
	return err
}

