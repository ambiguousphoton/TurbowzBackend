package main

import (
	"GoServer/repository"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

func recommendByVideo(w http.ResponseWriter, r *http.Request, vmdRepo repository.VMDrepo) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	videoIDStr := r.URL.Query().Get("video_id")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if videoIDStr == "" {
		http.Error(w, "video_id is required", http.StatusBadRequest)
		return
	}

	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid video_id", http.StatusBadRequest)
		return
	}

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

	results, err := vmdRepo.SimilaritySearch(videoID, limit, offset)
	if err != nil {
		http.Error(w, "Error during similarity search: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"page":    page,
		"limit":   limit,
		"results": results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("Sent similarity results for video_id %d (page %d)", videoID, page)
}

func recommendVideosByUserEmbedding(w http.ResponseWriter, r *http.Request, vmdRepo repository.VMDrepo) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

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

	results, err := vmdRepo.RecommendVideosByUserEmbedding(userID, limit, offset)
	if err != nil {
		http.Error(w, "Error during getting recommendations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"page":    page,
		"limit":   limit,
		"results": results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("Sent recommended video results for user_id %d (page %d)", userID, page)
}

func recommendEcosByUserEmbedding(w http.ResponseWriter, r *http.Request, ecoRepo repository.EcoRepo) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

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

	results, err := ecoRepo.RecommendEcosByUserEmbedding(userID, limit, offset)
	if err != nil {
		http.Error(w, "Error during getting recommendations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"page":    page,
		"limit":   limit,
		"results": results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("Sent recommended video results for user_id %d (page %d)", userID, page)
}





func main() {
	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)
	VMDrepo := repository.NewPostgresVMDRepo(db)
	ECOrepo := repository.NewPostgresEcoRepo(db)

	http.HandleFunc("/recommend-videos-for-user", func(w http.ResponseWriter, r *http.Request) {
		recommendVideosByUserEmbedding(w, r, VMDrepo)
	})

	http.HandleFunc("/recommend", func(w http.ResponseWriter, r *http.Request) {
		recommendByVideo(w, r, VMDrepo)
	})

	http.HandleFunc("/recommend-ecos", func(w http.ResponseWriter, r *http.Request) {
		recommendEcosByUserEmbedding(w, r, ECOrepo)
	})

	http.HandleFunc("/recommend-users", func(w http.ResponseWriter, r *http.Request) {
		recommendEcosByUserEmbedding(w, r, ECOrepo)
	})

	log.Println("Server started at port 8007")
	if err := http.ListenAndServe(":8007", nil); err != nil {
		log.Println("Critical Error:", err)
		os.Exit(1)
	}
}
