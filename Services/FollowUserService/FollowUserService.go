package main

import (
	"GoServer/authenticator"
	"GoServer/repository"
	"log"
	"net/http"
	"os"
	"strconv"
	_ "github.com/lib/pq"
	"encoding/json"
)


func FollowUser(w http.ResponseWriter, r *http.Request, userRepo repository.UserRepo) error{
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return nil
	}

	followerID := r.Context().Value("userID").(int64)
	
	followeeID, err := strconv.ParseInt(r.FormValue("followeeID"), 10, 64)
	if err != nil{
		http.Error(w, "Invalid Followee ID", http.StatusBadRequest)
		log.Println("Error parsing followee id to int64")
		return err
	}

	err = userRepo.FollowUser(followerID, followeeID)
	if err != nil{
		http.Error(w, "Failed to Follow User", http.StatusInternalServerError)
		log.Println("Error in userRepo.FollowUser ", err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	log.Println(followerID , " followed  ", followeeID)
	return nil
}

func UnfollowUser(w http.ResponseWriter, r *http.Request, userRepo repository.UserRepo) error{
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return nil
	}

	followerID := r.Context().Value("userID").(int64)
	
	followeeID, err := strconv.ParseInt(r.FormValue("followeeID"), 10, 64)
	if err != nil{
		http.Error(w, "Invalid Followee ID", http.StatusBadRequest)
		log.Println("Error parsing followee id to int64")
		return err
	}

	err = userRepo.UnfollowUser(followerID, followeeID)
	if err != nil{
		http.Error(w, "Failed to unfollow User", http.StatusInternalServerError)
		log.Println("Error in userRepo.UnfollowUser ", err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	log.Println(followerID , " unfollowed  ", followeeID)
	return nil
}



func getFollowers(w http.ResponseWriter, r *http.Request, userRepo repository.UserRepo) error{
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return nil
	}

	checkID,err :=  strconv.ParseInt(r.FormValue("checkID"), 10, 64)
	if err != nil{
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		log.Println("Error parsing followee id to int64", r.FormValue("checkID"))
		return err
	}

	followers, err := userRepo.GetAllFollowers(checkID)
	if err != nil{
		http.Error(w, "Failed to get followers", http.StatusInternalServerError)
		log.Println("Error in userRepo.GetFollowers ", err)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    response := map[string]interface{}{
        "followee_id": checkID,
        "followers":   followers,
    }

    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Println("Error encoding JSON:", err)
        return err
    }
	log.Println(checkID, " got followers  ", followers)
	return nil
}


func getFollowees(w http.ResponseWriter, r *http.Request, userRepo repository.UserRepo) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return nil
	}

	// Use query param for GET (?checkID=123)
	checkIDStr := r.URL.Query().Get("checkID")
	checkID, err := strconv.ParseInt(checkIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		log.Println("Error parsing follower id to int64:", checkIDStr)
		return err
	}

	followees, err := userRepo.GetAllFollowees(checkID)
	if err != nil {
		http.Error(w, "Failed to get followees", http.StatusInternalServerError)
		log.Println("Error in userRepo.GetFollowees:", err)
		return err
	}

	// Set headers before writing
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"GivenID":     checkID,
		"IsFollowing": followees,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding JSON:", err)
		return err
	}

	log.Println(checkID, "is following", followees)
	return nil
}



func main(){
	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	UserRepo := repository.NewPostgresUserRepo(db)
	
	http.HandleFunc("/follow", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request){
		err := FollowUser(w, r, UserRepo)
		if err != nil{
			http.Error(w, "Failed to Follow User" , http.StatusInternalServerError)
			return
		}
	}))

	http.HandleFunc("/unfollow", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request){
		err := UnfollowUser(w, r, UserRepo)
		if err != nil{
			http.Error(w, "Failed to Unfollow User" , http.StatusInternalServerError)
			return
		}
	}))


	http.HandleFunc("/get-followers", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request){
		err := getFollowers(w, r, UserRepo)
		if err != nil{
			http.Error(w, "failed get-followers" , http.StatusInternalServerError)
			return
		}
	}))

	http.HandleFunc("/get-followees", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request){
		err := getFollowees(w, r, UserRepo)
		if err != nil{
			http.Error(w, "Failed to get-followees" , http.StatusInternalServerError)
			return
		}
	}))
	
	log.Println("Server Started at Port 8010")
	err := http.ListenAndServe(":8010", nil)
	if err != nil{
		log.Println("Critical Error Occured", "error", err)
		os.Exit(1)
	}
}