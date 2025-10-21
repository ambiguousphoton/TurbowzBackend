package main

import (
	"GoServer/repository"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"strconv"
	"fmt"
)


func getVideos(w http.ResponseWriter, r *http.Request, vmdRepo repository.VMDrepo){
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return
	}

	keyword := r.URL.Query().Get("keyword")
    if keyword == "" {
        http.Error(w, "keyword is required", http.StatusBadRequest)
        return
    }

	results, err := vmdRepo.SearchVMDs(keyword)
	    if err != nil {
        http.Error(w,"Error while searching keyword Error:" + err.Error(), http.StatusInternalServerError)
        return
	}
	json.NewEncoder(w).Encode(results)
	log.Println("Search Results Sents")
}

func getVideosBy(w http.ResponseWriter, r *http.Request, vmdRepo repository.VMDrepo)error{
	log.Println("getVideosByUserID called")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return fmt.Errorf("wrong metod %v", http.StatusMethodNotAllowed)
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

	results, err := vmdRepo.SearchVMDsBy(userID)
	    if err != nil {
        http.Error(w,"Error while searching video with userID Error:" + err.Error(), http.StatusInternalServerError)
        return fmt.Errorf("error while searching video with userID: %w", err)
	}
	json.NewEncoder(w).Encode(results)
	log.Println("video Search Results Sents for userID:", userID)
	return nil
}

func main(){
	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	VMDrepo := repository.NewPostgresVMDRepo(db)

	log.Println("Server Started at Port 8082") 
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
    	getVideos(w, r, VMDrepo)
	})

	http.HandleFunc("/search-video-with", func(w http.ResponseWriter, r *http.Request) {
    	err := getVideosBy(w, r, VMDrepo)
		if err != nil{
			log.Println("Error in getVideosByUserID:", err)
		}
	})


	err := http.ListenAndServe(":8082", nil)
	if err != nil{
		log.Println("Critical Error Occured", "error", err)
		os.Exit(1)
	}

}