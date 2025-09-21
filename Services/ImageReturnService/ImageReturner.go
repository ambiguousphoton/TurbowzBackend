package main

import (
	"log"
	"net/http"
	"os"
)


func getImage(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Content-Type", "image/jpeg")
	
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return
	}

	img_url := r.URL.Query().Get("img")
	if img_url == ""{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("No Image URL")
		return
	}

	if img_url[0] == 'p'{
		img_path := "./MediaData/UserProfile/" +  img_url + ".jpg"
		http.ServeFile(w, r, img_path)
		log.Println("img served")
		return
	}

	img_path := "./MediaData/thumbnails/" +  img_url + ".jpg"
	http.ServeFile(w, r, img_path)
	log.Println("img served")
}





func main(){
	log.Println("Server Started at Port 8088") 
	http.HandleFunc("/i", func(w http.ResponseWriter, r *http.Request) {
    	getImage(w, r)
	})

	err := http.ListenAndServe(":8088", nil)
	if err != nil{
		log.Println("Critical Error Occured", "error", err)
		os.Exit(1)
	}

}