package repository

import (
        "database/sql"
        "log"
        "GoServer/models"
		"strings"
		"fmt"
        )

type UserRepo interface {
    CreateNewUser(user *models.UserData, auth *models.UserAuth) error
	GetUser(userID int64) (*models.UserData, error)
	UpadateUserProfile(user *models.UserData) error
    AddUserAuth(auth *models.UserAuth) error
	CheckUser(userHandle string) (int64, string, error)
	FollowUser(FollowerID int64, FolloweeID int64) (error)
	UnfollowUser(FollowerID int64, FolloweeID int64) (error)
	GetAllFollowers(FolloweeID int64) ([]int64,error)
	GetAllFollowees(followerID int64) ([]int64 ,error)
	AddConnection(requesterID int64, respondentID int64) (error)
	SearchWithKeyword(keyword string) ([]int64, error)
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
		INSERT INTO user_data_table (user_handle, user_profile_name, user_description, from_location,  gender)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING user_id
	`
	err = tx.QueryRow(userQuery,
		user.UserHandle,
		user.UserProfileName,
		user.UserDescription,
		user.FromLocation,
		// user.UserDateOfBirth,
		user.Gender,
	).Scan(&user.UserID)
	if err != nil {
        log.Printf("Error inserting into user_data_table: %v", err)
        return err
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
		log.Printf("Error inserting into user_authentication table")
        return err
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
