package main

import (
	"GoServer/repository"
	"log"
	"net/http"
	"os"
	"encoding/json"
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





func main(){
	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	VMDrepo := repository.NewPostgresVMDRepo(db)

	log.Println("Server Started at Port 8082") 
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
    	getVideos(w, r, VMDrepo)
	})

	err := http.ListenAndServe(":8082", nil)
	if err != nil{
		log.Println("Critical Error Occured", "error", err)
		os.Exit(1)
	}

}