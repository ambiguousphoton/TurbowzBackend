package main

import (
	"GoServer/authenticator"
	"GoServer/models"
	"GoServer/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"        // postgres driver
	"golang.org/x/crypto/bcrypt" // encrytion
)


var jwtKey = []byte("om namo bhagwate vaudevay")  /// SUPER SECRET KEEP SOME WHERE SAFE
var tokenValidityDuration = time.Minute * 50000




func createNewUser(w http.ResponseWriter, r * http.Request, UserRepo repository.UserRepo) error{
	// Only Post method is allowed
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("createNewUser: Method not allowed - received %s, expected POST", r.Method)
		return nil
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("createNewUser: Failed to parse form data - %v", err)
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return err
	}
	

	userHandle := r.FormValue("user_handle")
	profileName := r.FormValue("user_profile_name")
	userDescription := r.FormValue("userDescription")
	fromLocation := r.FormValue("fromLocation")
	userGender := r.FormValue("gender")
	email := r.FormValue("email")
	phoneNumber := r.FormValue("phoneNumber")
	password := r.FormValue("password")

	// Validate all required fields
	if err := authenticator.ValidateUserHandle(userHandle); err != nil {
		log.Printf("createNewUser: validation failed - %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	if err := authenticator.ValidateProfileName(profileName); err != nil {
		log.Printf("createNewUser: validation failed - %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	if err := authenticator.ValidateEmail(email); err != nil {
		log.Printf("createNewUser: validation failed for handle %s - %v", userHandle, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	if err := authenticator.ValidatePhone(phoneNumber); err != nil {
		log.Printf("createNewUser: validation failed for handle %s - %v", userHandle, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	if err := authenticator.ValidatePassword(password); err != nil {
		log.Printf("createNewUser: password validation failed for handle %s", userHandle)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	url := uuid.New().String()
	userhashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("createNewUser: Failed to hash password for user %s - %v", userHandle, err)
		http.Error(w, "Problem with Password", http.StatusInternalServerError)
		return err
	}
	
	user := &models.UserData{
		UserHandle:      userHandle,
		UserProfileName: profileName,
		UserDescription: userDescription,
		FromLocation:    fromLocation,
		// UserDateOfBirth: userDateOfBirth,
		Gender:          userGender,
		Url: 			 url,	

	}

	auth := &models.UserAuth{
		UserLoginAccount:   email,
		UserPhoneNumber:    phoneNumber,
		UserHashedPassword: string(userhashedPassword),
	}

	err = UserRepo.CreateNewUser(user, auth)
	if err != nil {
		log.Printf("createNewUser: Failed to create user in database for handle %s - %v", userHandle, err)
		return err
	}
	log.Printf("User created with UserHandle = %v, UserID = %v, AuthID = %v\n",user.UserHandle, user.UserID, auth.AuthID)

	// Create JWT
	expirationTime := time.Now().Add(tokenValidityDuration)
	claims := &models.Claims{
		UserID:     user.UserID,
		UserHandle: user.UserHandle,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("createNewUser: Failed to sign JWT token for user %d - %v", user.UserID, err)
		http.Error(w, "could not create token", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
    "token": tokenString,
	})
	log.Println("User Authenticated with id", user.UserID)
	return nil
}




func upadateProfile(w http.ResponseWriter, r * http.Request, UserRepo repository.UserRepo) error{
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("updateProfile: Method not allowed - received %s, expected POST", r.Method)
		return nil
	}
	err := r.ParseForm()
	if err != nil {
		log.Printf("updateProfile: Failed to parse form data - %v", err)
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return err
	}

	userIDint, ok := r.Context().Value("userID").(int64)
	if !ok  {
		log.Printf("updateProfile: Invalid or missing userID in context")
		http.Error(w, "Error InvalidUserId ", http.StatusBadRequest)
		return fmt.Errorf("InvalidUserId")
	}

	newUser := &models.UserData{
		UserID          : userIDint,
		UserHandle      : r.FormValue("user_handle"),
		UserProfileName : r.FormValue("user_profile_name"), 
		UserDescription : r.FormValue("userDescription"),
		FromLocation    : r.FormValue("fromLocation"),
		Gender          : r.FormValue("gender"),
	}

	err = UserRepo.UpadateUserProfile(newUser)
	
	if err!= nil{
		log.Printf("updateProfile: Failed to update profile for user %d - %v", userIDint, err)
		http.Error(w, "Failed to update user profile: "+err.Error(), http.StatusInternalServerError)
		return err
	}


	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User profile updated successfully"))

	log.Println("Information Updated, for user with id: ", userIDint)
	return nil
}






func UserAuthentication(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo) error{
	if r.Method != http.MethodPost {
		log.Printf("UserAuthentication: Method not allowed - received %s, expected POST", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return fmt.Errorf("wrong method: %s", r.Method)
	}
	err := r.ParseForm()
	if err != nil {
		log.Printf("UserAuthentication: Failed to parse form data - %v", err)
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return err
	}

	UserHandle := r.FormValue("user_handle")
	inputPassword := r.FormValue("password")

	if UserHandle == "" || inputPassword == "" {
		log.Printf("UserAuthentication: missing user_handle or password")
		http.Error(w, "user_handle and password are required", http.StatusBadRequest)
		return fmt.Errorf("missing credentials")
	}

	userId, PasswordFrmDB, err := UserRepo.CheckUser(UserHandle)
	if err != nil{
		log.Printf("UserAuthentication: Failed to find user with handle %s - %v", UserHandle, err)
		http.Error(w, "Failed to Authenticate User, UserHandle not found", http.StatusInternalServerError)
		return err
	} 
	if bcrypt.CompareHashAndPassword([]byte(PasswordFrmDB), []byte(inputPassword)) != nil{
		log.Printf("UserAuthentication: Password mismatch for user %s (ID: %d)", UserHandle, userId)
		http.Error(w, "Failed to Authenticate User, Wrong Password", http.StatusInternalServerError)
		return err
	}

		// Create JWT
	expirationTime := time.Now().Add(tokenValidityDuration)
	claims := &models.Claims{
		UserID:     userId,
		UserHandle: UserHandle,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("UserAuthentication: Failed to sign JWT token for user %d - %v", userId, err)
		http.Error(w, "could not create token", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
    "token": tokenString,
	"userID": fmt.Sprintf("%d", userId),
	})
	log.Println("User Authenticated with id", userId)
	return nil
}


func GetUser(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo) error{
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("GetUser: Method not allowed - received %s, expected GET", r.Method)
		return nil
	}

	userIDstr := r.URL.Query().Get("userID")
	if userIDstr == ""{
		http.Error(w, "userID is required", http.StatusBadRequest)
		log.Println("GetUser: userID query parameter is missing")
		return fmt.Errorf("userID is required")
	}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil{
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		log.Printf("GetUser: Invalid userID format for input %s - %v", userIDstr, err)
		return fmt.Errorf("invalid userID format")
	}
	
	userData, err := UserRepo.GetUser(userID)
	if err != nil{
		log.Printf("GetUser: Failed to retrieve user data for userID %s - %v", userIDstr, err)
		http.Error(w, "Error while fetching user data: "+err.Error(), http.StatusInternalServerError)
		return err
	}
	json.NewEncoder(w).Encode(userData)
	log.Println("Search Results for userID Sent :", userID)
	return nil
}


func SearchUsers(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo) error{
	log.Println("SearchUsers: SearchUsers endpoint hit")
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("SearchUsers: Method not allowed - received %s, expected GET", r.Method)
		return nil
	}
	keyword := r.URL.Query().Get("keyword")
	if keyword == ""{
		http.Error(w, "keyword is required", http.StatusBadRequest)
		log.Println("keyword: parameter is missing")
		return fmt.Errorf("keyword is required")
	}



	log.Println("SearchUsers: Searching users with keyword:", keyword)
	
	userIDList, err := UserRepo.SearchWithKeyword(keyword)
	if err != nil{
		log.Printf("SearchUsers: Failed to retrieve user list for keyword %s - %v", keyword, err)
		http.Error(w, "Error while fetching user list : "+err.Error(), http.StatusInternalServerError)
		return err
	}
	json.NewEncoder(w).Encode(userIDList)
	log.Println("Search Results for keyword Sent")
	return nil
}

func 	saveVideoHandler(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo) error{
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("SaveVideo: Method not allowed - received %s, expected POST", r.Method)
		return nil
	}

	userID, ok := r.Context().Value("userID").(int64)
	if !ok  {
		log.Printf("saveVideoHandler: Invalid or missing userID in context")
		http.Error(w, "Error InvalidUserId ", http.StatusBadRequest)
		return fmt.Errorf("InvalidUserId")
	}

	videoIDStr := r.URL.Query().Get("videoID")
	if videoIDStr == ""{
		http.Error(w, "videoID is required", http.StatusBadRequest)
		log.Println("saveVideoHandler: videoID query parameter is missing")
		return 	fmt.Errorf("videoID is required")
	}

	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil{
		http.Error(w, "Invalid videoID format", http.StatusBadRequest)
		log.Printf("saveVideoHandler: Invalid videoID format for input %s - %v", videoIDStr, err)
		return fmt.Errorf("invalid videoID format")
	}

	saved, err := UserRepo.UserSavedVideo(userID, videoID)
	if err != nil{
		http.Error(w, "Error saving watched video: "+err.Error(), http.StatusInternalServerError)
		log.Printf("saveVideoHandler: Error saving video for user_id %d and video_id %d - %v", userID, videoID, err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"saved": saved})
	log.Printf("Saved watched video for user_id %d and video_id %d", userID, videoID)
	return nil
}

func saveEcoHandler(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo) error{
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("SaveEco: Method not allowed - received %s, expected POST", r.Method)
		return nil
	}

	userID, ok := r.Context().Value("userID").(int64)
	if !ok  {
		log.Printf("saveEcoHandler: Invalid or missing userID in context")
		http.Error(w, "Error InvalidUserId ", http.StatusBadRequest)
		return fmt.Errorf("InvalidUserId")
	}

	ecoIDstr := r.URL.Query().Get("ecoID")
	ecoID, err := strconv.ParseInt(ecoIDstr, 10, 64)
	if err != nil{
		http.Error(w, "Invalid ecoID format", http.StatusBadRequest)
		log.Printf("saveEcoHandler: Invalid ecoID format for input %s - %v", ecoIDstr, err)
		return fmt.Errorf("invalid ecoID format")
	}


	saved,err := UserRepo.UserSavedEco(userID, ecoID)
	if err != nil{
		http.Error(w, "Error checking saved status for ECO data: "+err.Error(), http.StatusInternalServerError)
		log.Printf("saveEcoHandler: Error checking saved status for for user_id %d - %v", userID, err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"saved": saved})
	log.Printf("Saved ECO data for user_id %d", userID)
	return nil
}

func getEcoSavedStatus(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo) error{
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("getEcoSavedStatus: Method not allowed - received %s, expected POST", r.Method)
		return nil
	}

	userID, ok := r.Context().Value("userID").(int64)
	if !ok  {
		log.Printf("getEcoSavedStatus: Invalid or missing userID in context")
		http.Error(w, "Error InvalidUserId ", http.StatusBadRequest)
		return fmt.Errorf("InvalidUserId")
	}

	ecoIDstr := r.URL.Query().Get("ecoID")
	ecoID, err := strconv.ParseInt(ecoIDstr, 10, 64)
	if err != nil{
		http.Error(w, "Invalid ecoID format", http.StatusBadRequest)
		log.Printf("getEcoSavedStatus: Invalid ecoID format for input %s - %v", ecoIDstr, err)
		return fmt.Errorf("invalid ecoID format")
	}


	saved,err := UserRepo.UserEcoSavedStatus(userID, ecoID)
	if err != nil{
		http.Error(w, "Error saving ECO data: "+err.Error(), http.StatusInternalServerError)
		log.Printf("getEcoSavedStatus: Error saving ECO data for user_id %d - %v", userID, err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"saved": saved})
	log.Printf("getEcoSavedStatus ECO data for user_id %d", userID)
	return nil
}

func getVideoSavedStatus(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo) error{
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("getVideoSavedStatus: Method not allowed - received %s, expected POST", r.Method)
		return nil
	}

	userID, ok := r.Context().Value("userID").(int64)
	if !ok  {
		log.Printf("getVideoSavedStatus: Invalid or missing userID in context")
		http.Error(w, "Error InvalidUserId ", http.StatusBadRequest)
		return fmt.Errorf("InvalidUserId")
	}

	videoIDStr := r.URL.Query().Get("videoID")
	if videoIDStr == ""{
		http.Error(w, "videoID is required", http.StatusBadRequest)
		log.Println("getVideoSavedStatus: videoID query parameter is missing")
		return 	fmt.Errorf("videoID is required")
	}

	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil{
		http.Error(w, "Invalid videoID format", http.StatusBadRequest)
		log.Printf("getVideoSavedStatus: Invalid videoID format for input %s - %v", videoIDStr, err)
		return fmt.Errorf("invalid videoID format")
	}

	saved, err := UserRepo.UserVideoSavedStatus(userID, videoID)
	if err != nil{
		http.Error(w, "Error saving watched video: "+err.Error(), http.StatusInternalServerError)
		log.Printf("getVideoSavedStatus: Error checking saved status for  video for user_id %d and video_id %d - %v", userID, videoID, err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"saved": saved})
	log.Printf("Saved watched video for user_id %d and video_id %d", userID, videoID)
	return nil
}


func turbomaxStatusCheck(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo) {
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("turbomaxStatusCheck: Method not allowed - received %s, expected GET", r.Method)
		return
	}
	userIDStr := r.URL.Query().Get("userID")
	if userIDStr == ""{
		http.Error(w, "userID is required", http.StatusBadRequest)
		log.Println("turbomaxStatusCheck: userID query parameter is missing")
		return 	
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil{
		http.Error(w, "Invalid userIDStr format", http.StatusBadRequest)
		log.Printf("turbomaxStatusCheck: Invalid userID format for input %s - %v", userIDStr, err)
		return 
	}
	status, err := UserRepo.GetTurbomaxStatusOfUser(userID)
	if err != nil{
		http.Error(w,"Error Occured while getting status", http.StatusInternalServerError)
		log.Println("Error Occured while getting status from db : ",err.Error())
	}
	log.Println("Returning status of turbomax for userID", userID, "status", status)
	json.NewEncoder(w).Encode(map[string]bool{"turbomax_active": status})
}


func verifyEmail(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo, ev *authenticator.EmailVerifier) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Printf("verifyEmail: failed to parse form - %v", err)
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	if err := authenticator.ValidateEmail(email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exists, err := UserRepo.CheckEmailExists(email)
	if err != nil {
		log.Printf("verifyEmail: db error checking email - %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if exists {
		log.Printf("verifyEmail: email already registered - %s", email)
		http.Error(w, "email is already registered", http.StatusConflict)
		return
	}

	if err := ev.GenerateAndSend(email); err != nil {
		log.Printf("verifyEmail: failed to send OTP to %s - %v", email, err)
		http.Error(w, "failed to send verification email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "verification code sent"})
}

func confirmEmail(w http.ResponseWriter, r *http.Request, ev *authenticator.EmailVerifier) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Printf("confirmEmail: failed to parse form - %v", err)
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	code := r.FormValue("otp")
	if email == "" || code == "" {
		http.Error(w, "email and otp are required", http.StatusBadRequest)
		return
	}

	if !ev.Verify(email, code) {
		log.Printf("confirmEmail: invalid or expired OTP for %s", email)
		http.Error(w, "invalid or expired verification code", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "email verified"})
}

func forgotPassword(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo, ev *authenticator.EmailVerifier) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	if err := authenticator.ValidateEmail(email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exists, err := UserRepo.CheckEmailExists(email)
	if err != nil {
		log.Printf("forgotPassword: db error - %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !exists {
		log.Printf("forgotPassword: email not found - %s", email)
		http.Error(w, "email not registered", http.StatusNotFound)
		return
	}

	if err := ev.GenerateAndSend(email); err != nil {
		log.Printf("forgotPassword: failed to send OTP to %s - %v", email, err)
		http.Error(w, "failed to send verification email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "password reset code sent"})
}

func resetPassword(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo, ev *authenticator.EmailVerifier) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	code := r.FormValue("otp")
	newPassword := r.FormValue("new_password")

	if email == "" || code == "" || newPassword == "" {
		http.Error(w, "email, otp, and new_password are required", http.StatusBadRequest)
		return
	}

	if err := authenticator.ValidatePassword(newPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !ev.Verify(email, code) {
		log.Printf("resetPassword: invalid or expired OTP for %s", email)
		http.Error(w, "invalid or expired verification code", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("resetPassword: failed to hash password - %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := UserRepo.UpdatePassword(email, string(hashedPassword)); err != nil {
		log.Printf("resetPassword: failed to update password for %s - %v", email, err)
		http.Error(w, "failed to reset password", http.StatusInternalServerError)
		return
	}

	log.Printf("resetPassword: password reset successful for %s", email)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "password reset successful"})
}

func withCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}


func main() {
	
	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	UserRepo := repository.NewPostgresUserRepo(db)

	// Email verification — SMTP config
	emailVerifier := authenticator.NewEmailVerifier(
		"smtp.gmail.com",
		"587",
		"turbowz.official@gmail.com",
		"pijr hppk ycqy rkxo",
		"turbowz.official@gmail.com",
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/verify-email", func(w http.ResponseWriter, r *http.Request) {
		verifyEmail(w, r, UserRepo, emailVerifier)
	})

	mux.HandleFunc("/confirm-email", func(w http.ResponseWriter, r *http.Request) {
		confirmEmail(w, r, emailVerifier)
	})

	mux.HandleFunc("/forgot-password", func(w http.ResponseWriter, r *http.Request) {
		forgotPassword(w, r, UserRepo, emailVerifier)
	})

	mux.HandleFunc("/reset-password", func(w http.ResponseWriter, r *http.Request) {
		resetPassword(w, r, UserRepo, emailVerifier)
	})

	mux.HandleFunc("/create-new-account", func(w http.ResponseWriter, r *http.Request) {
    	err := createNewUser(w, r, UserRepo)
		if err != nil {
			log.Printf("Handler /create-new-account: %v", err)
			http.Error(w, "Failed to create new user", http.StatusInternalServerError)
			return
		}
	})


	mux.HandleFunc("/update-profile", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request){
		err := upadateProfile(w, r, UserRepo)
		if err != nil{
			log.Printf("Handler /update-profile: %v", err)
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}
	}))


	mux.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request){
		err := UserAuthentication(w, r, UserRepo)
		if err != nil{
			log.Printf("Handler /authenticate: %v", err)
			http.Error(w, "Failed to Authenticate User" , http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/get-user", func(w http.ResponseWriter, r *http.Request){
		err := GetUser(w, r, UserRepo)
		if err != nil{
			log.Printf("Handler /get-user: %v", err)
			http.Error(w, "Failed to get User" , http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/search-users", func(w http.ResponseWriter, r *http.Request){
		err := SearchUsers(w, r, UserRepo)
		if err != nil{
			log.Printf("Handler /search-users: %v", err)
			http.Error(w, "Failed to search for Users" , http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/video-saved-status", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		getVideoSavedStatus(w, r, UserRepo)
	}))

	mux.HandleFunc("/eco-saved-status", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		getEcoSavedStatus(w, r, UserRepo)
	}))

	mux.HandleFunc("/save-video", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		saveVideoHandler(w, r, UserRepo)
	}))

	mux.HandleFunc("/save-eco", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {	
		
		saveEcoHandler(w, r, UserRepo)
	}))


	mux.HandleFunc("/get-turbomax-status", func(w http.ResponseWriter, r*http.Request){
		turbomaxStatusCheck(w,r, UserRepo)
	})
	
	log.Println("UserDataService Started at Port 8100")
	err := http.ListenAndServe(":8100", withCORS(mux))
	if err != nil{
		log.Printf("Critical Error: Failed to start server on port 8100 - %v", err)
		os.Exit(1)
	}
}

