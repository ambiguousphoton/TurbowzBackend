package main

import (
	"GoServer/Services/ServerDataReceive/tools"
	"GoServer/authenticator"
	"GoServer/models"
	"GoServer/repository"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"image/jpeg"
	"image"
	"github.com/google/uuid"
	"fmt"
)


func saveVideoOnServer(w http.ResponseWriter, r *http.Request, vmdRepo repository.VMDrepo) {


    if r.Method != http.MethodPost{
		log.Printf("saveVideoOnServer: Method not allowed - received %s, expected POST", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, "Method not Allowed", http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		log.Printf("saveVideoOnServer: Failed to parse multipart form - %v", err)
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}


	//// ----------------------------  Initialising directory

	mediaDataDir := "./MediaData/videos/rawVideos"        
	if err := os.MkdirAll(mediaDataDir, os.ModePerm); err != nil {
		log.Printf("saveVideoOnServer: Failed to create directory %s - %v", mediaDataDir, err)
		http.Error(w, "Unable to create media data directory", http.StatusInternalServerError)
		return
	}

	////------------------------------ ------------------------------------------------------


	////------------------------------ Saving Video to rawVideo Folder

	video_url := uuid.New().String()

	receivedFile, fileHeader, err := r.FormFile("video")
	if err != nil {
		log.Printf("saveVideoOnServer: Failed to retrieve video file from form - %v", err)
		http.Error(w, "Error retrieving file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer receivedFile.Close()



	// Extract original file extension
	videoExtenstion := filepath.Ext(fileHeader.Filename) // e.g., ".mp4", ".avi"

	mediaBinaryPath := filepath.Join(mediaDataDir , video_url + videoExtenstion)
	binaryFile, err := os.Create(mediaBinaryPath)
		if err != nil {
		log.Printf("saveVideoOnServer: Failed to create file %s - %v", mediaBinaryPath, err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer binaryFile.Close()



	if _, err := io.Copy(binaryFile, receivedFile); err != nil {
		log.Printf("saveVideoOnServer: Failed to write video data to file %s - %v", mediaBinaryPath, err)
		http.Error(w, "Unable to write file", http.StatusInternalServerError)
		return
	}


	//// --------------------------  Saving Meta Data ----------------------------------


	uploader_id, ok := r.Context().Value("userID").(int64)
	if !ok  {
		log.Printf("saveVideoOnServer: Invalid or missing userID in context")
		http.Error(w, "Error Invalid UserId ", http.StatusBadRequest)
		return 
	}


	title := r.FormValue("title")
	info := r.FormValue("info")
	user_name := r.FormValue("user_name")
	tagsJSON := r.FormValue("tags")
	

	// extract tags from JSON array
	var tags []string
	if tagsJSON != "" {
		err := json.Unmarshal([]byte(tagsJSON), &tags)
		if err != nil {
			http.Error(w, "Invalid JSON in 'tags' field", http.StatusBadRequest)
			return
		}
	}

	vmd := &models.VideoMetaData{
		Uploader_ID:     uploader_id, 
		Title:           title,
		Video_Info:      info,
		Video_Url:       video_url,
		Tags: 		 	 tags,
	}
	video_id, err := vmdRepo.CreateNewVMD(vmd)
	if err != nil {
		log.Printf("saveVideoOnServer: Failed to save video metadata for user %d, video %s - %v", uploader_id, video_url, err)
		http.Error(w, "Unable to save metadata: "+err.Error(), http.StatusInternalServerError)
		return
    }


	
	///-------------------------------------- end Save meta data to DB

	/////// -------------------------- Compression & Encoding -----------------------------------------

	binaryFile.Sync()

	tools.EncodeVideoToHLS(video_url)

	if err = tools.FirstFrameThumbnail(video_url); err != nil{
		log.Printf("saveVideoOnServer: Failed to create thumbnail for video %s - %v", video_url, err)
	}

	// tools.EncodeVideoToDASH(uid)

	if tags == nil {
			tags = []string{}
		}
	// Call vectorization API
	vectorData := map[string]interface{}{
		"title": title,
		"description": info,
		"tags": tags ,
		"user_name": user_name,
		"video_id": video_id,
	}
	
	jsonData, _ := json.Marshal(vectorData)
	resp, err := http.Post("http://localhost:9000/vectorize-video/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("saveVideoOnServer: Failed to call vectorization API - %v", err)
	} else {
		resp.Body.Close()
		log.Printf("saveVideoOnServer: Video vectorization completed for video %s", video_url)
	}
	
	log.Printf("saveVideoOnServer: Successfully uploaded video %s for user %d (title: %s)", video_url, uploader_id, title)

}


func saveECO(w http.ResponseWriter, r *http.Request, postRepo repository.EcoRepo){
	if r.Method != http.MethodPost{
		log.Printf("makePost: Method not allowed - received %s, expected POST", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, "Method not Allowed", http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(10 << 21)
	if err != nil {
		log.Printf("Failed to parse multipart form - %v", err)
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}
	eco_url := uuid.New().String()
	eco_text := r.FormValue("eco_text")
	uploader_name := r.FormValue("uploader_name")
	uploader_id, ok := r.Context().Value("userID").(int64)
	
	if !ok  {
		log.Printf("makePost: Invalid or missing userID in context")
		http.Error(w, "Error Invalid UserId ", http.StatusBadRequest)
		return 
	}
	
	tagsJSON := r.FormValue("tags")
	var tags []string
	if tagsJSON != "" {
		err := json.Unmarshal([]byte(tagsJSON), &tags)
		if err != nil {
			http.Error(w, "Invalid JSON in 'tags' field", http.StatusBadRequest)
			return
		}
	}

	uploadDir := "./MediaData/EcoImages/"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	files := r.MultipartForm.File["images"];
	var imagePaths []string;
	for imageIndex, fileHeader := range files {
		log.Printf("Processing uploaded image: %d", imageIndex)
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("Failed to open uploaded file: %v", err)
			http.Error(w, "Error processing file: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			log.Printf("Failed to decode image %s: %v", fileHeader.Filename, err)
			return
		}
		eco_image_uuid := eco_url + "_" + strconv.Itoa(imageIndex)
		eco_image_filepath := filepath.Join(uploadDir, eco_image_uuid + ".jpg")
		savedFile, err := os.Create(eco_image_filepath)
		if err != nil {
			log.Printf("Failed to create image file %s: %v", eco_image_filepath, err)
			http.Error(w, "Error saving image: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer savedFile.Close()
		
		err = jpeg.Encode(savedFile, img, &jpeg.Options{Quality: 90})
		if err != nil {
			log.Printf("Failed to encode JPEG for %s: %v", eco_image_filepath, err)
			http.Error(w, "Error encoding image: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		savedFile.Sync()
		imagePaths = append(imagePaths, eco_image_filepath)
	}
	EcoPost := &models.EcoMetaData{
		Uploader_ID:  		uploader_id,
		Eco_Text:     		eco_text,
		Tags:				tags,
		Images_Count:		len(imagePaths),
		Eco_Url: 			eco_url,
	}
	eco_id, err := postRepo.CreateEcoPost(EcoPost)
	if err != nil {
		log.Printf("Failed to save ECO post metadata for user %d, ECO %s - %v", uploader_id, eco_url, err)
		http.Error(w, "Unable to save metadata: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	log.Printf("Successfully created ECO post %s with ID %d for user %d", eco_url, eco_id, uploader_id)
	vectorData := map[string]interface{}{
		"eco_text": eco_text,
		"tags": tags ,
		"img_count": len(imagePaths),
		"uploader_name": uploader_name,
		"eco_id": eco_id,
	}
	jsonData, _ := json.Marshal(vectorData)
	resp, err := http.Post("http://localhost:9000/vectorize-eco/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("saveECO: Failed to call vectorization API - %v", err)
	} else {
		resp.Body.Close()
		log.Printf("saveECO: Eco vectorization completed for video %s", eco_url)
	}

	log.Printf("saveECO: Successfully uploaded ECO %s for user %d", eco_url, uploader_id)
}

func saveEvent(w http.ResponseWriter, r *http.Request, eventRepo repository.EventRepo) {
	if r.Method != http.MethodPost {
		log.Printf("saveEvent: Method not allowed - received %s, expected POST", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, "Method not Allowed", http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(10 << 21)
	if err != nil {
		log.Printf("Failed to parse multipart form - %v", err)
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	event_url := uuid.New().String()
	uploader_id, ok := r.Context().Value("userID").(int64)
	if !ok {
		log.Printf("saveEvent: Invalid or missing userID in context")
		http.Error(w, "Error Invalid UserId", http.StatusBadRequest)
		return
	}

	event_title := r.FormValue("event_title")
	event_description := r.FormValue("event_description")
	event_start_time := r.FormValue("event_start_time")
	event_end_time := r.FormValue("event_end_time")

	tagsJSON := r.FormValue("tags")
	var tags []string
	if tagsJSON != "" {
		err := json.Unmarshal([]byte(tagsJSON), &tags)
		if err != nil {
			http.Error(w, "Invalid JSON in 'tags' field", http.StatusBadRequest)
			return
		}
	}

	uploadDir := "./MediaData/EventImages/"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["images"]
	var imagePaths []string
	for imageIndex, fileHeader := range files {
		log.Printf("Processing uploaded image: %d", imageIndex)
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("Failed to open uploaded file: %v", err)
			http.Error(w, "Error processing file: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			log.Printf("Failed to decode image %s: %v", fileHeader.Filename, err)
			return
		}

		event_image_uuid := event_url + "_" + strconv.Itoa(imageIndex)
		event_image_filepath := filepath.Join(uploadDir, event_image_uuid+".jpg")
		savedFile, err := os.Create(event_image_filepath)
		if err != nil {
			log.Printf("Failed to create image file %s: %v", event_image_filepath, err)
			http.Error(w, "Error saving image: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer savedFile.Close()

		err = jpeg.Encode(savedFile, img, &jpeg.Options{Quality: 90})
		if err != nil {
			log.Printf("Failed to encode JPEG for %s: %v", event_image_filepath, err)
			http.Error(w, "Error encoding image: "+err.Error(), http.StatusInternalServerError)
			return
		}

		savedFile.Sync()
		imagePaths = append(imagePaths, event_image_filepath)
	}

	eventData := &models.EventMetaData{
		Uploader_ID:      uploader_id,
		Event_Url:        event_url,
		Event_Title:      event_title,
		Event_Description: event_description,
		Tags:             tags,
		Images_Count:     len(imagePaths),
		Event_Start_Time: event_start_time,
		Event_End_Time:   event_end_time,
	}

	event_id, err := eventRepo.CreateNewEvent(eventData)
	if err != nil {
		log.Printf("Failed to save event metadata for user %d, event %s - %v", uploader_id, event_url, err)
		http.Error(w, "Unable to save metadata: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("Successfully created event %s with ID %d for user %d", event_url, event_id, uploader_id)
}

func savePFP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("savePFP: Method not allowed - received %s, expected POST", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit upload size to 10 MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Failed to parse multipart form - %v", err)
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// ✅ Fix 1: Type assert userID properly
	uploaderID, ok := r.Context().Value("userID").(int64)
	if !ok {
		log.Printf("savePFP: Invalid or missing userID in context")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// ✅ Fix 2: Get the uploaded file correctly
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("Failed to retrieve profile image: %v", err)
		http.Error(w, "Missing profile image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create user profile folder if it doesn’t exist
	uploadDir := "./MediaData/UserProfile/"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		log.Printf("Failed to decode uploaded image %s: %v", header.Filename, err)
		http.Error(w, "Invalid image file", http.StatusBadRequest)
		return
	}

	// ✅ Fix 3: filepath.Join arguments must be strings — convert uploaderID
	savePath := filepath.Join(uploadDir, fmt.Sprintf("%d.jpg", uploaderID))

	outFile, err := os.Create(savePath)
	if err != nil {
		log.Printf("Failed to create file %s: %v", savePath, err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})
	if err != nil {
		log.Printf("Failed to encode JPEG for %s: %v", savePath, err)
		http.Error(w, "Error encoding image", http.StatusInternalServerError)
		return
	}

	outFile.Sync()

	w.WriteHeader(http.StatusCreated)
	log.Printf("Successfully saved profile picture for user %d ", uploaderID)
	
}



func main(){

	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	VMDrepo := repository.NewPostgresVMDRepo(db)
	EcoRepo := repository.NewPostgresEcoRepo(db)
	EventRepo := repository.NewPostgresEventRepo(db)



	log.Println("Server Started at Port 8080")
	http.HandleFunc("/upload", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
    	saveVideoOnServer(w, r, VMDrepo)
	}))

	http.HandleFunc("/eco-upload", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
    	saveECO(w, r, EcoRepo)
	}))

	http.HandleFunc("/pfp-upload", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
    	savePFP(w, r)
	}))
	
	http.HandleFunc("/event-upload", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
    	saveEvent(w, r, EventRepo)
	}))

	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		log.Printf("Critical Error: Failed to start server on port 8080 - %v", err)
		os.Exit(1)
	}
}