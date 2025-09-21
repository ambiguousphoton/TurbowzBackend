package main

import (
	"GoServer/authenticator"
	"GoServer/models"
	"GoServer/repository"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"encoding/json"
)

func pushComment(w http.ResponseWriter, r *http.Request, CommentRepo repository.CommentRepo) error {
	if r.Method != http.MethodPost{
		http.Error(w, "Method not Allowed", http.StatusBadRequest)
		log.Println("method not allowed")
		return fmt.Errorf( "method not Allowed")
	}
	// Parse the request body
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		log.Println("failed to parse request body")
		return err
	}
	
	commmenter_id, ok := r.Context().Value("userID").(int64)
	if !ok  {
		http.Error(w, "Error Invalid UserId ", http.StatusBadRequest)
		log.Println("error Invalid UserId ", commmenter_id)
		return  fmt.Errorf("error Invalid UserId ")
	}

	parent_video_id, err := strconv.ParseInt(r.FormValue("parentVideoID"), 10, 64)
	if err != nil {
		http.Error(w, "Error Invalid Video ID ", http.StatusBadRequest)
		log.Println("error in string to int conversion ", parent_video_id)
		return err
	}
	comment_text := r.FormValue("commentText")

	new_comment := &models.CommentData{
		Commenter_id: commmenter_id,
		Parent_video_id: parent_video_id,
		Comment_text: comment_text,
	}

	cmnt_ID, err :=CommentRepo.CreateNewComment(new_comment)

	if err != nil{
		http.Error(w, "Error in creating comment", http.StatusBadRequest)
		return fmt.Errorf("error from CommentDataRepo.go in CreatingNewComment %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Commented"))
	log.Println("comment added with: id ", cmnt_ID, " on video: ", parent_video_id)

	return nil
}


func getComments(w http.ResponseWriter, r *http.Request, CommentRepo repository.CommentRepo) error{
	if r.Method != http.MethodGet{
		http.Error(w, "Method not Allowed", http.StatusBadRequest)
		log.Println("method not allowed")
		return fmt.Errorf( "method not Allowed")
	}
	video_id_str := r.URL.Query().Get("videoID")
	video_id, err := strconv.ParseInt(video_id_str, 10, 64)
	if err != nil{
		http.Error(w, "Error Invalid Video ID ", http.StatusBadRequest)
		log.Println("error in string to int conversion ", video_id)
		return err
	}

	comments, err := CommentRepo.GetVideoComments(video_id)
	if err != nil{
		http.Error(w, "Error in getting comments", http.StatusBadRequest)
		log.Println("error in getting comments ", err)
		return err
	}
	json.NewEncoder(w).Encode(comments)
	log.Println("comments fetched for video: ", video_id)
	return nil
}


func main(){

	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	CommentRepo := repository.NewPostgresCommentRepo(db)

	log.Println("Server Started at Port 7200")
	http.HandleFunc("/push-comment", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
    	err := pushComment(w, r, CommentRepo)
		if err != nil{
			log.Print(err)
		}
	}))

	http.HandleFunc("/get-comment", func(w http.ResponseWriter, r *http.Request) {
    	err := getComments(w, r, CommentRepo)
		if err != nil{
			log.Print(err)
		}
	})
	
	err := http.ListenAndServe(":7200", nil)
	if err != nil{
		log.Println("Critical Error Occured", "error", err)
		os.Exit(1)
	}
}