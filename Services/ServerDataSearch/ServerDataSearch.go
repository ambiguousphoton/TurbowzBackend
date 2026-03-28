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

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	
	limit := 20
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

	results, err := vmdRepo.SearchVMDs(keyword, limit, offset)
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

func getEcoByUserId(w http.ResponseWriter, r *http.Request, ecoRepo repository.EcoRepo)error{
	log.Println("getEcoByUserId called")
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
		return err
	}

	results, err := ecoRepo.SearchEcosByUserID(userID)
	    if err != nil {
        http.Error(w,"Error while searching video with userID Error:" + err.Error(), http.StatusInternalServerError)
    	log.Printf("getEcoByUserId: error while searching video with userID: %v", err)
		return err
	}
	json.NewEncoder(w).Encode(results)
	log.Println("getEcoByUserId: eco of userID: ", userID, " sent successfully")
	return nil
}


func main(){
	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	VMDrepo := repository.NewPostgresVMDRepo(db)
	ECOrepo := repository.NewPostgresEcoRepo(db)

	log.Println("Server Started at Port 8082") 
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
    	getVideos(w, r, VMDrepo)
	})

	http.HandleFunc("/search-video-with", func(w http.ResponseWriter, r *http.Request) {
    	err := getVideosBy(w, r, VMDrepo)
		if err != nil{
			log.Println("Failed in getVideosByUserID:", err)
		}
	})

	http.HandleFunc("/search-eco-by-user", func(w http.ResponseWriter, r *http.Request) {
    	err := getEcoByUserId(w, r, ECOrepo)
			if err != nil{
				log.Println("Failed in getEcoByUserId:", err)
		}
	})

	err := http.ListenAndServe(":8082", nil)
	if err != nil{
		log.Println("Critical Error Occured", "error", err)
		os.Exit(1)
	}

}