package main

import (
	"GoServer/repository"
	"log"
	"net/http"
	"os"
	"fmt"
	"encoding/json"
	"strconv"
	_ "github.com/lib/pq" // postgres driver
	"GoServer/authenticator"
)



func getUserWatchHistory(w http.ResponseWriter, r *http.Request, VMDRepo repository.VMDrepo) error{
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("GetUser: Method not allowed - received %s, expected GET", r.Method)
		return nil
	}

	userID, ok := r.Context().Value("userID").(int64)
	if !ok  {
		log.Printf("getUserWatchHistory: Invalid or missing userID in context")
		http.Error(w, "Error InvalidUserId ", http.StatusBadRequest)
		return fmt.Errorf("InvalidUserId")
	}
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	// default pagination values
	page := 1
	limit := 10
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	offset := (page - 1) * limit

	WatchedHistory, err := VMDRepo.GetUserWatchHistory(userID, limit, offset)
	if err != nil {
		http.Error(w, "Error during similarity search: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	response := map[string]interface{}{
		"page":    page,
		"limit":   limit,
		"results": WatchedHistory,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("Sent Watch history results for user_id %d (page %d)", userID, page)
	return nil
}

func getSecuredUserWatchHistory(w http.ResponseWriter, r *http.Request, VMDRepo repository.VMDrepo) error{
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("GetUser: Method not allowed - received %s, expected GET", r.Method)
		return nil
	}

	
	userIDstr := r.URL.Query().Get("userID")
	if userIDstr == ""{
		http.Error(w, "userID is required", http.StatusBadRequest)
		log.Println("GetUser: userID query parameter is missing")
		return 	fmt.Errorf("userID is required")
	}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil{
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		log.Printf("GetUser: Invalid userID format for input %s - %v", userIDstr, err)
		return fmt.Errorf("invalid userID format")
	}

	password := r.URL.Query().Get("password")
	if password != "OMAngOMAngOMang"{
		return fmt.Errorf("password not correct")
	}
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	// default pagination values
	page := 1
	limit := 10
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	offset := (page - 1) * limit

	WatchedHistory, err := VMDRepo.GetUserWatchHistory(userID, limit, offset)
	if err != nil {
		http.Error(w, "Error during similarity search: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	response := map[string]interface{}{
		"page":    page,
		"limit":   limit,
		"results": WatchedHistory,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("Sent Watch history results for user_id %d (page %d)", userID, page)
	return nil
}


func deleteMyHistory(w http.ResponseWriter, r *http.Request, VMDRepo repository.VMDrepo) error{
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("GetUser: Method not allowed - received %s, expected GET", r.Method)
		return nil
	}
	userIDstr := r.URL.Query().Get("userID")
	if userIDstr == ""{
		http.Error(w, "userID is required", http.StatusBadRequest)
		log.Println("deleteMyHistory: userID query parameter is missing")
		return 	fmt.Errorf("userID is required")
	}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil{
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		log.Printf("deleteMyHistory: Invalid userID format for input %s - %v", userIDstr, err)
		return fmt.Errorf("invalid userID format")
	}


	err = VMDRepo.DeleteMyHistory(userID)
	if err != nil {
		http.Error(w, "Error during deleteMyHistory : "+err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("History Deleted")
	log.Printf("History Deleted for user_id %d ", userID)
	return nil
}

func getUserActivityData(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("getUserActivityData: Method not allowed - received %s, expected GET", r.Method)
		return nil
	}
	userIDstr := r.URL.Query().Get("userID")
	if userIDstr == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		log.Println("getUserActivityData: userID query parameter is missing")
		return fmt.Errorf("userID is required")
	}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		log.Printf("getUserActivityData: Invalid userID format for input %s - %v", userIDstr, err)
		return fmt.Errorf("invalid userID format")
	}

	analyticsData, err := UserRepo.GetUserAnalytics(userID)
	if err != nil {
		http.Error(w, "Error during getUserActivityData: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analyticsData)
	log.Printf("Sent activity data for user_id %d", userID)
	return nil
}

func postVideoVote(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("postVideoVote: Method not allowed - received %s, expected POST", r.Method)
		return nil
	}

	var req struct {
		VideoID int64 `json:"video_id"`
		UserID  int64 `json:"user_id"`
		Quality int   `json:"quality"`
		AiUsage int   `json:"ai_usage"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("postVideoVote: JSON decode error - %v", err)
		return err
	}

	// basic validation
	if req.VideoID == 0 || req.UserID == 0 {
		http.Error(w, "video_id and user_id are required", http.StatusBadRequest)
		return fmt.Errorf("missing required fields")
	}

	if req.Quality < 1 || req.Quality > 5 || req.AiUsage < 1 || req.AiUsage > 5 {
		http.Error(w, "quality and ai_usage must be between 1 and 5", http.StatusBadRequest)
		return fmt.Errorf("invalid vote range")
	}

	err := UserRepo.UpsertVideoVote(
		req.VideoID,
		req.UserID,
		req.Quality,
		req.AiUsage,
	)
	if err != nil {
		http.Error(w, "Error posting video vote: "+err.Error(), http.StatusInternalServerError)
		log.Printf("postVideoVote: repo error - %v", err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "video vote recorded",
	})
	log.Printf("Video vote recorded: video_id=%d user_id=%d", req.VideoID, req.UserID)

	return nil
}

func postEchoVote(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("postEchoVote: Method not allowed - received %s, expected POST", r.Method)
		return nil
	}

	var req struct {
		EcoID   int64 `json:"eco_id"`
		UserID  int64 `json:"user_id"`
		Quality int   `json:"quality"`
		AiUsage int   `json:"ai_usage"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("postEchoVote: JSON decode error - %v", err)
		return err
	}

	if req.EcoID == 0 || req.UserID == 0 {
		http.Error(w, "eco_id and user_id are required", http.StatusBadRequest)
		return fmt.Errorf("missing required fields")
	}

	if req.Quality < 1 || req.Quality > 5 || req.AiUsage < 1 || req.AiUsage > 5 {
		http.Error(w, "quality and ai_usage must be between 1 and 5", http.StatusBadRequest)
		return fmt.Errorf("invalid vote range")
	}

	err := UserRepo.UpsertEcoVote(
		req.EcoID,
		req.UserID,
		req.Quality,
		req.AiUsage,
	)
	if err != nil {
		http.Error(w, "Error posting eco vote: "+err.Error(), http.StatusInternalServerError)
		log.Printf("postEchoVote: repo error - %v", err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "eco vote recorded",
	})
	log.Printf("Eco vote recorded: eco_id=%d user_id=%d", req.EcoID, req.UserID)

	return nil
}

func main() {
	
	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	VMDrepo := repository.NewPostgresVMDRepo(db)
	UserRepo := repository.NewPostgresUserRepo(db)
    http.HandleFunc("/get-user-watch-history", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		getUserWatchHistory(w, r, VMDrepo)
	}))

	http.HandleFunc("/secured-get-user-watch-history", func(w http.ResponseWriter, r *http.Request) {
		getSecuredUserWatchHistory(w, r, VMDrepo)
	})

	http.HandleFunc("/delete-my-history", func(w http.ResponseWriter, r *http.Request) {
		deleteMyHistory(w, r, VMDrepo)
	})

	http.HandleFunc("/get-activity-data", func(w http.ResponseWriter, r *http.Request){
		getUserActivityData(w, r, UserRepo)
	})

	http.HandleFunc("/post-video-vote", func(w http.ResponseWriter, r *http.Request){
		postVideoVote(w, r, UserRepo)
	})

	http.HandleFunc("/post-echo-vote", func(w http.ResponseWriter, r *http.Request){
		postEchoVote(w, r, UserRepo)
	})

	log.Println("Activity Server Started at Port 7992")
	err := http.ListenAndServe(":7992", nil)
	if err != nil{
		log.Println("Critical Error Occured", "error", err)
		os.Exit(1)
	}
}