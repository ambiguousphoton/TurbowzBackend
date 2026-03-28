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


	img_path := "./MediaData/thumbnails/" +  img_url + ".jpg"
	http.ServeFile(w, r, img_path)
	log.Println("img served")
}


func getEcoImage(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Content-Type", "image/jpeg")
	
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return
	}

	img_url := r.URL.Query().Get("eco_url")
	index   := r.URL.Query().Get("index")
	if img_url == ""{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("No Image URL")
		return
	}



	img_path := "./MediaData/EcoImages/" +  img_url +"_"+ index + ".jpg"
	http.ServeFile(w, r, img_path)
	log.Println("eco img served")
}

func getProfileImage(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Content-Type", "image/jpeg")
	
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return
	}

	user_id := r.URL.Query().Get("user_id")
	if user_id == ""{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("No Image URL")
		return
	}


	
	img_path := "./MediaData/UserProfile/" +  user_id  +".jpg"
	log.Println(img_path)
	http.ServeFile(w, r, img_path)

	log.Println("pfp img served")
}

func getBannerAd(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Content-Type", "image/jpeg")
	
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return
	}

	ad_id := r.URL.Query().Get("ad_id")
	if ad_id == ""{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("No Ad ID")
		return
	}
	img_path := "./MediaData/BannerAds/" +  ad_id  +".jpg"
	log.Println(img_path)
	http.ServeFile(w, r, img_path)

	log.Println("ad img served")
}

func getEventImage(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Content-Type", "image/jpeg")
	
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return
	}

	event_url := r.URL.Query().Get("event_url")
	if event_url == ""{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("No Event URL")
		return
	}
	img_path := "./MediaData/EventImages/" +  event_url  +".jpg"
	log.Println(img_path)
	http.ServeFile(w, r, img_path)

	log.Println("event img served")
}

func main(){
	log.Println("Server Started at Port 8088") 
	http.HandleFunc("/i", func(w http.ResponseWriter, r *http.Request) {
    	getImage(w, r)
	})

	http.HandleFunc("/e", func(w http.ResponseWriter, r *http.Request) {
    	getEcoImage(w, r)
	})
	http.HandleFunc("/pfp", func(w http.ResponseWriter, r *http.Request) {
    	getProfileImage(w, r)
	})

	http.HandleFunc("/ad", func(w http.ResponseWriter, r *http.Request) {
    	getBannerAd(w, r)
	})

	http.HandleFunc("/event-img", func(w http.ResponseWriter, r *http.Request) {
    	getEventImage(w, r)
	})
	err := http.ListenAndServe(":8088", nil)
	if err != nil{
		log.Println("Critical Error Occured", "error", err)
		os.Exit(1)
	}

}

