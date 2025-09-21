package main

import (
	"GoServer/authenticator"
	"GoServer/repository"
	"log"
	"net/http"
	"os"
	"strconv"
)


func  addConnection(w http.ResponseWriter, r *http.Request, UserRepo repository.UserRepo) error{
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return nil
	}

	var requesterID int64 = 27
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

	log.Println("Server Started at Port 8001")
	err := http.ListenAndServe(":8001", nil)
	if err != nil{
		log.Println("Critical Error Occured", "error", err)
		os.Exit(1)
	}
}