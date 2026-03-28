package main

import (

	"GoServer/authenticator"
	"GoServer/models"
	"GoServer/repository"

	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"image/jpeg"
	"image"
	"encoding/json"
	"fmt"
)


func uploadBannerAd(w http.ResponseWriter, r *http.Request, adsRepo repository.AdsRepo) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Extract uploader ID
	uploaderID, ok := r.Context().Value("userID").(int64)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	redirectURL := r.FormValue("redirect_url")

	if title == "" || redirectURL == "" {
		http.Error(w, "title and redirect_url are required", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Missing image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Invalid image", http.StatusBadRequest)
		return
	}

	// Insert into DB BEFORE saving file
	adData := &models.BannerAd{
		UploaderID:  uploaderID,
		Title:       title,
		RedirectURL: redirectURL,
	}

	adID, err := adsRepo.CreateNewBannerAd(adData)
	if err != nil {
		http.Error(w, "DB insert error", http.StatusInternalServerError)
		log.Printf("DB insert error: %v", err)
		return
	}

	// Save image as <ad_id>.jpg
	uploadDir := "./MediaData/BannerAds/"
	os.MkdirAll(uploadDir, os.ModePerm)

	savePath := filepath.Join(uploadDir, fmt.Sprintf("%d.jpg", adID))

	outFile, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Error saving image", http.StatusInternalServerError)
		log.Printf("Error creating file: %v", err)
		return
	}
	defer outFile.Close()

	jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"success","ad_id":` + strconv.FormatInt(adID, 10) + `}`))
}

func getBannerAds(w http.ResponseWriter, r *http.Request, adsRepo repository.AdsRepo) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	ads, err := adsRepo.GetBannerAds(page, limit)
	if err != nil {
		http.Error(w, "Failed to fetch ads", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ads)
}

func main(){

	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	adsRepo := repository.NewPostgresAdsRepo(db)

	log.Println("Server Started at Port 8991")


	http.HandleFunc("/upload-b-ads",authenticator.RequireAuth( func(w http.ResponseWriter, r *http.Request) {
    	uploadBannerAd(w, r, adsRepo)
	}))
	
	http.HandleFunc("/get-b-ads", func(w http.ResponseWriter, r *http.Request) {
	getBannerAds(w, r, adsRepo)
	})


	err := http.ListenAndServe(":8991", nil)
	if err != nil{
		log.Printf("Critical Error: Failed to start server on port 8080 - %v", err)
		os.Exit(1)
	}
}