package main

import(
	"net/http"
	"os"
	"log"
	"path/filepath"
	"github.com/google/uuid"
	"io"
	"GoServer/repository"
	"GoServer/models"
	"GoServer/Services/ServerDataReceive/tools"
	"GoServer/authenticator"
	
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


	vmd := &models.VideoMetaData{
		Uploader_ID:     uploader_id, 
		Title:           title,
		Video_Info:      info,
		Video_Url:       video_url,
	}

	if err := vmdRepo.CreateNewVMD(vmd); err != nil {
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

	/////// -------------------------- End Compression & Encoding -----------------------------------------





	log.Printf("saveVideoOnServer: Successfully uploaded video %s for user %d (title: %s)", video_url, uploader_id, title)

}


func main(){

	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	VMDrepo := repository.NewPostgresVMDRepo(db)

	log.Println("Server Started at Port 8080")
	http.HandleFunc("/upload", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
    	saveVideoOnServer(w, r, VMDrepo)
	}))


	
	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		log.Printf("Critical Error: Failed to start server on port 8080 - %v", err)
		os.Exit(1)
	}
}