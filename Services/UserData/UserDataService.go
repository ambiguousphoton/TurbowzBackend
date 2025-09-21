package main

import (
	"GoServer/models"
	"GoServer/repository"
	"log"
	"net/http"
	"os"
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"GoServer/authenticator"
	_ "github.com/lib/pq" // postgres driver
	"golang.org/x/crypto/bcrypt" // encrytion
	"encoding/json"
)


var jwtKey = []byte("om namo bhagwate vaudevay")  /// SUPER SECRET KEEP SOME WHERE SAFE





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
	// userDateOfBirth := r.FormValue("userDateOfBirth")
	userGender := r.FormValue("gender")
	email := r.FormValue("email")	
	phoneNumber := r.FormValue("phoneNumber")
	password := r.FormValue("password")

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
	expirationTime := time.Now().Add(15 * time.Minute)
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
	expirationTime := time.Now().Add(15 * time.Minute)
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
	})
	log.Println("User Authenticated with id", userId)
	return nil
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

	mux := http.NewServeMux()

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



	
	log.Println("Server Started at Port 8100")
	err := http.ListenAndServe(":8100", withCORS(mux))
	if err != nil{
		log.Printf("Critical Error: Failed to start server on port 8100 - %v", err)
		os.Exit(1)
	}
}
