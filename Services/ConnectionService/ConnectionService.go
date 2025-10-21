package main

import (
	"GoServer/authenticator"
	"GoServer/repository"
	"log"
	"net/http"
	"os"
	"strconv"
	"fmt"
)


func  addConnection(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo) error{
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return nil
	}

	requesterID, ok := r.Context().Value("userID").(int64)
	if !ok  {
		http.Error(w, "Error Invalid UserId ", http.StatusBadRequest)
		log.Println("error Invalid UserId ", requesterID)
		return  fmt.Errorf("error Invalid UserId ")
	}

	respondentID, err := strconv.ParseInt(r.FormValue("contactID"), 10, 64)
	if err != nil{
		return err
	}
	err = UserRepo.AddConnection(requesterID, respondentID)
	if err != nil{
		log.Println("failed to connect users")
		return err
	}
	w.WriteHeader(http.StatusCreated)
	log.Println("Connection added in ", requesterID, respondentID)
	return err
}



func getUsersConnections(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo) error{
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return nil
	}
	return nil
}



func main(){
	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	UserRepo := repository.NewPostgresUserRepo(db)

	http.HandleFunc("/add-connection", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request){
		err := addConnection(w, r, UserRepo)
		if err != nil{
			log.Println(err)
			return
		}
	}))

	http.HandleFunc("/get-user-connection", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request){
		err := addConnection(w, r, UserRepo)
		if err != nil{
			log.Println(err)
			return
		}
	}))


	log.Println("Server Started at Port 8001")
	err := http.ListenAndServe(":8001", nil)
	if err != nil{
		log.Println("Critical Error Occured", "error", err)
		os.Exit(1)
	}
}