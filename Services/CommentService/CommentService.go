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
	"database/sql"
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

	var parentCommentID sql.NullInt64

	parentCommentStr := r.FormValue("parentCommentID")

	if parentCommentStr != "" {
		val, err := strconv.ParseInt(parentCommentStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid parent comment ID", http.StatusBadRequest)
			log.Println("string to int conversion error:", err)
			return err
		}
		parentCommentID = sql.NullInt64{
			Int64: val,
			Valid: true,
		}
	} else {
		parentCommentID = sql.NullInt64{Valid: false} // NULL
	}

	comment_text := r.FormValue("commentText")

	log.Printf("Creating comment - VideoID: %d, ParentCommentID: %v, CommenterID: %d", parent_video_id, parentCommentID, commmenter_id)

	new_comment := &models.CommentData{
		Commenter_id: commmenter_id,
		Parent_video_id: parent_video_id,
		Comment_text: comment_text,
		Parent_Comment_ID: parentCommentID,
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

	limit := 20
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	comments, err := CommentRepo.GetVideoComments(video_id, limit, offset)
	if err != nil{
		http.Error(w, "Error in getting comments", http.StatusBadRequest)
		log.Println("error in getting comments ", err)
		return err
	}

	total, err := CommentRepo.GetVideoCommentsCount(video_id)
	if err != nil{
		log.Println("error getting comment count", err)
		total = 0
	}

	response := map[string]interface{}{
		"comments": comments,
		"pagination": map[string]interface{}{
			"limit": limit,
			"offset": offset,
			"total": total,
		},
	}

	json.NewEncoder(w).Encode(response)
	log.Println("comments fetched for video: ", video_id)
	return nil
}

func pushEcoComment(w http.ResponseWriter, r *http.Request, CommentRepo repository.CommentRepo) error {
	if r.Method != http.MethodPost{
		http.Error(w, "Method not Allowed", http.StatusBadRequest)
		log.Println("method not allowed")
		return fmt.Errorf( "method not Allowed")
	}

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

	parent_eco_id, err := strconv.ParseInt(r.FormValue("parentEcoID"), 10, 64)
	if err != nil {
		http.Error(w, "Error Invalid Eco ID ", http.StatusBadRequest)
		log.Println("error in string to int conversion ", parent_eco_id)
		return err
	}

	var parentCommentID sql.NullInt64

	parentCommentStr := r.FormValue("parentCommentID")

	if parentCommentStr != "" {
		val, err := strconv.ParseInt(parentCommentStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid parent comment ID", http.StatusBadRequest)
			log.Println("string to int conversion error:", err)
			return err
		}
		parentCommentID = sql.NullInt64{
			Int64: val,
			Valid: true,
		}
	} else {
		parentCommentID = sql.NullInt64{Valid: false}
	}

	comment_text := r.FormValue("commentText")

	log.Printf("Creating eco comment - EcoID: %d, ParentCommentID: %v, CommenterID: %d", parent_eco_id, parentCommentID, commmenter_id)

	new_comment := &models.EcoCommentData{
		Commenter_id: commmenter_id,
		Parent_Eco_id: parent_eco_id,
		Comment_text: comment_text,
		Parent_Comment_ID: parentCommentID,
	}

	cmnt_ID, err :=CommentRepo.CreateNewEcoComment(new_comment)

	if err != nil{
		http.Error(w, "Error in creating eco comment", http.StatusBadRequest)
		return fmt.Errorf("error from CommentDataRepo.go in CreatingNewEcoComment %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Eco Commented"))
	log.Println("eco comment added with: id ", cmnt_ID, " on eco: ", parent_eco_id)

	return nil
}

func getEcoComments(w http.ResponseWriter, r *http.Request, CommentRepo repository.CommentRepo) error{
	if r.Method != http.MethodGet{
		http.Error(w, "Method not Allowed", http.StatusBadRequest)
		log.Println("method not allowed")
		return fmt.Errorf( "method not Allowed")
	}
	eco_id_str := r.URL.Query().Get("ecoID")
	eco_id, err := strconv.ParseInt(eco_id_str, 10, 64)
	if err != nil{
		http.Error(w, "Error Invalid Eco ID ", http.StatusBadRequest)
		log.Println("error in string to int conversion ", eco_id)
		return err
	}

	limit := 20
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	comments, err := CommentRepo.GetEcoComments(eco_id, limit, offset)
	if err != nil{
		http.Error(w, "Error in getting eco comments", http.StatusBadRequest)
		log.Println("error in getting eco comments ", err)
		return err
	}

	total, err := CommentRepo.GetEcoCommentsCount(eco_id)
	if err != nil{
		log.Println("error getting eco comment count", err)
		total = 0
	}

	response := map[string]interface{}{
		"comments": comments,
		"pagination": map[string]interface{}{
			"limit": limit,
			"offset": offset,
			"total": total,
		},
	}

	json.NewEncoder(w).Encode(response)
	log.Println("eco comments fetched for eco: ", eco_id)
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

	http.HandleFunc("/push-eco-comment", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
    	err := pushEcoComment(w, r, CommentRepo)
		if err != nil{
			log.Print(err)
		}
	}))

	http.HandleFunc("/get-eco-comment", func(w http.ResponseWriter, r *http.Request) {
    	err := getEcoComments(w, r, CommentRepo)
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