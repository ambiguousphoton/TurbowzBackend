package main

import(
	"net/http"
	"strings"
	"path/filepath"
	"log"
)


func getVideoStream(w http.ResponseWriter, r * http.Request){
	if r.Method != http.MethodGet {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
	log.Println("Narayan Narayan request received.", r.URL.Path)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	videoDestination := "./MediaData/videos/hlsEncodedVideos" + strings.TrimPrefix(r.URL.Path, "/get-video-stream" )
	ext := filepath.Ext(videoDestination)
	log.Println(ext, videoDestination)

	switch ext {
		case ".m3u8":
			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		case ".ts":
			w.Header().Set("Content-Type", "video/mp2t")
		default:
			w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeFile(w, r, videoDestination)
}



func main(){
	http.HandleFunc("/get-video-stream/", getVideoStream)
	log.Println("Streaming Server started at http://localhost:8091")
	err := http.ListenAndServe(":8091", nil)
	if err != nil{
		log.Println(" Error starting server: ", err)
		return 
	}
}