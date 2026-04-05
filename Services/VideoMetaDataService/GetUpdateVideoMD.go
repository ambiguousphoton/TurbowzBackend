package main

import (
	"GoServer/authenticator"
	// "GoServer/models"
	"GoServer/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq" // postgres driver
)


func getVideoMD(w http.ResponseWriter, r *http.Request, vmdRepo repository.VMDrepo) error{
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return fmt.Errorf("wrong Method")
	}
    if err := r.ParseForm(); err != nil {
        log.Printf("Failed to parse form: %v\n", err)
        http.Error(w, "failed to parse form", http.StatusBadRequest)
        return  fmt.Errorf("failed to parse form")
    }
	userID, ok := r.Context().Value("userID").(int64)
	if !ok  {
		log.Printf("updateProfile: Invalid or missing userID in context")
		http.Error(w, "Error InvalidUserId ", http.StatusBadRequest)
		return  fmt.Errorf("wrong InvalidUserId")
	}


	video_id_str := r.URL.Query().Get("video_id")
	video_id, err := strconv.ParseInt(video_id_str, 10, 64)
	if err != nil{
        http.Error(w, "invalid video_id", http.StatusBadRequest)
		return fmt.Errorf("invalid video_idr")
	}


	results, err := vmdRepo.GetSpecificVideoMD(video_id, userID)
	if err != nil {
		http.Error(w,"Error while searching videoID Error:" + err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("error while searching videoID Error: %s" , err.Error())
	}
	json.NewEncoder(w).Encode(results)
	log.Println("Search Results Sents", video_id)
	return nil
}



func videoViewUpdate(w http.ResponseWriter, r *http.Request, vmdRepo repository.VMDrepo, userRepo repository.UserRepo) {
    log.Println("videoViewUpdate called") 

    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method != http.MethodPost {
        log.Printf("Wrong Method: %s\n", r.Method)
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    if err := r.ParseForm(); err != nil {
        log.Printf("Failed to parse form: %v\n", err)
        http.Error(w, "failed to parse form", http.StatusBadRequest)
        return
    }

    videoIDStr := r.FormValue("video_id")
    userIDStr  := r.FormValue("user_id")
    log.Printf("Received video_id: %s\n", videoIDStr)

    videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
    if err != nil {
        log.Printf("Invalid video_id: %s, error: %v\n", videoIDStr, err)
        http.Error(w, "invalid video_id", http.StatusBadRequest)
        return
    }
    if userIDStr == "" {
        log.Printf("Not able to record History, User not signed In, no userID")
    }else {
        userID, err := strconv.ParseInt(userIDStr, 10, 64)
        if err != nil {
            log.Printf("Invalid userID: %s, error: %v\n", userIDStr, err)
            http.Error(w, "invalid userID", http.StatusBadRequest)
        } else{
            err = userRepo.AddVideoInUserHistory(userID, videoID)
            if err != nil{
                log.Printf("Error Updating UserHistory %v", err)
            }
            log.Printf("Video add in User History")
        }
    }

    err = vmdRepo.VideoViewUpdate(videoID)
    if err != nil {
        log.Printf("Error updating views for video_id %d: %v\n", videoID, err)
        http.Error(w, "failed to update views", http.StatusInternalServerError)
        return
    }

    log.Printf("Views updated successfully for video_id %d\n", videoID)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("views updated"))
}


func updateLuv(w http.ResponseWriter, r *http.Request, vmdRepo repository.VMDrepo) {
    log.Println("videoViewUpdate called") 

    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method != http.MethodPost {
        log.Printf("Wrong Method: %s\n", r.Method)
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    if err := r.ParseForm(); err != nil {
        log.Printf("Failed to parse form: %v\n", err)
        http.Error(w, "failed to parse form", http.StatusBadRequest)
        return
    }
	userID, ok := r.Context().Value("userID").(int64)
	if !ok  {
		log.Printf("updateProfile: Invalid or missing userID in context")
		http.Error(w, "Error InvalidUserId ", http.StatusBadRequest)
		return 
	}


    videoIDStr := r.FormValue("video_id")
    log.Printf("Received video_id: %s\n", videoIDStr)

    videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
    if err != nil {
        log.Printf("Invalid video_id: %s, error: %v\n", videoIDStr, err)
        http.Error(w, "invalid video_id", http.StatusBadRequest)
        return
    }

    luved, err := vmdRepo.UpdateLuv(videoID, userID)
    if err != nil {
        log.Printf("Error updating luvs for video_id %d and user_id %d: %v\n", videoID, userID, err)
        http.Error(w, "failed to update luv", http.StatusInternalServerError)
        return
    }

    log.Printf("luv updated successfully for video_id %d\n", videoID)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]bool{"luved": luved})
}

func getSavedVideos(w http.ResponseWriter, r *http.Request, vmdRepo repository.VMDrepo) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	videos, err := vmdRepo.GetSavedVideos(userID, limit, offset)
	if err != nil {
		log.Printf("Error fetching saved videos for user_id %d: %v", userID, err)
		http.Error(w, "Failed to get saved videos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

func getTrendingVideos(w http.ResponseWriter, r *http.Request, vmdRepo repository.VMDrepo) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userIDstr := r.URL.Query().Get("userID")
	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil{
		http.Error(w, "invalid userID", http.StatusBadRequest)
		return
	}
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	videos, err := vmdRepo.GetTrendingVMDsPaginated(userID, limit, offset)
	if err != nil {
		log.Printf("Error fetching trending videos %d: %v", err)
		http.Error(w, "Failed to get trending videos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

func getVideoScore(w http.ResponseWriter, r *http.Request, vmdRepo repository.VMDrepo){
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return 
	}
    if err := r.ParseForm(); err != nil {
        log.Printf("Failed to parse form: %v\n", err)
        http.Error(w, "failed to parse form", http.StatusBadRequest)
        return  
    }

	video_id_str := r.URL.Query().Get("video_id")
	video_id, err := strconv.ParseInt(video_id_str, 10, 64)
	if err != nil{
        http.Error(w, "invalid video_id", http.StatusBadRequest)
		return
	}

	video_score, err := vmdRepo.GetVideoScore(video_id)
	if err != nil {
        log.Printf("Error Getting Score for video_id %d: %v\n", video_id, err)
        http.Error(w, "failed to score", http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(video_score)	
	if err != nil{
		log.Printf("Error Sending the Vidoe Score for video ID %d, getting error %v", video_id, err)
	}
	log.Printf("Successfuly sent the Video Score for video ID %d", video_id)
}


func main() {
	
	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	VMDrepo := repository.NewPostgresVMDRepo(db)
    UserRepo :=  repository.NewPostgresUserRepo(db)

	http.HandleFunc("/vmd", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
    	err := getVideoMD(w, r, VMDrepo)
		if err != nil {
			http.Error(w, "Failed to get Video Meta Data", http.StatusInternalServerError)
			return
		}
	}))

	http.HandleFunc("/luv", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
    	updateLuv(w, r, VMDrepo)
	}))

	http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
    	videoViewUpdate(w, r, VMDrepo, UserRepo)	
	})

	http.HandleFunc("/get-saved-videos", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		getSavedVideos(w, r, VMDrepo)
	}))

	http.HandleFunc("/get-trending-videos", func(w http.ResponseWriter, r *http.Request) {
		getTrendingVideos(w, r, VMDrepo)
	})


	http.HandleFunc("/get-videos-score", func(w http.ResponseWriter, r *http.Request) {
		getVideoScore(w, r, VMDrepo)
	})

	log.Println("VMD GETER Server Started at Port 7999")
	err := http.ListenAndServe(":7999", nil)
	if err != nil{
		log.Println("Critical Error Occured", "error", err)
		os.Exit(1)
	}
}
