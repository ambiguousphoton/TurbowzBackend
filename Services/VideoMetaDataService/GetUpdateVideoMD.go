package main

import (
	"GoServer/repository"
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
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

	video_id_str := r.URL.Query().Get("video_id")
	video_id, err := strconv.ParseInt(video_id_str, 10, 64)
	if err != nil{
        http.Error(w, "invalid video_id", http.StatusBadRequest)
		return fmt.Errorf("invalid video_idr")
	}


	results, err := vmdRepo.GetSpecificVideoMD(video_id)
	if err != nil {
		http.Error(w,"Error while searching videoID Error:" + err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("error while searching videoID Error: %s" , err.Error())
	}
	json.NewEncoder(w).Encode(results)
	log.Println("Search Results Sents", video_id)
	return nil
}



func videoViewUpdate(w http.ResponseWriter, r *http.Request, vmdRepo repository.VMDrepo) {
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
    log.Printf("Received video_id: %s\n", videoIDStr)

    videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
    if err != nil {
        log.Printf("Invalid video_id: %s, error: %v\n", videoIDStr, err)
        http.Error(w, "invalid video_id", http.StatusBadRequest)
        return
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



func main() {
	
	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	VMDrepo := repository.NewPostgresVMDRepo(db)


	http.HandleFunc("/vmd", func(w http.ResponseWriter, r *http.Request) {
    	err := getVideoMD(w, r, VMDrepo)
		if err != nil {
			http.Error(w, "Failed to get Video Meta Data", http.StatusInternalServerError)
			return
		}
	})


	http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
    	videoViewUpdate(w, r, VMDrepo)	
	})

	log.Println("VMD GETER Server Started at Port 7999")
	err := http.ListenAndServe(":7999", nil)
	if err != nil{
		log.Println("Critical Error Occured", "error", err)
		os.Exit(1)
	}
}
