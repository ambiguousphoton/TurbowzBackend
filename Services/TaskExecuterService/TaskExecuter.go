package main

import (

	"GoServer/repository"
	"fmt"
	"log"
	"net/http"
	"os"

)

func updateUserEmbeddings(w http.ResponseWriter, r *http.Request, userRepo repository.UserRepo, vmdRepo repository.VMDrepo) error {
	userIDParam := r.URL.Query().Get("userID")
	password := r.URL.Query().Get("OMAngOMAngOMang")

	// --- If updating all users, require password ---
	if userIDParam == "*" {
		if password != "OMAngOMAngOMang" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return fmt.Errorf("unauthorized: password required for global update")
		}

		const batchSize = 100
		const historyCount = 20
		offset := 0
		totalUpdated := 0

		log.Println("🔸 Updating embeddings for all users in batches...")

		for {
			userIDs, err := userRepo.AllUsersReturn(batchSize, offset)
			if err != nil {
				log.Printf("❌ Failed to fetch users (offset %d): %v", offset, err)
				return err
			}

			if len(userIDs) == 0 {
				log.Printf("✅ All users processed. Total updated: %d", totalUpdated)
				break
			}

			for _, id := range userIDs {
				log.Printf("➡️ Updating embeddings for user_id: %d", id)

				err := vmdRepo.UpdateUserEmbeddingsFromVideoHistory(id, historyCount)
				if err != nil {
					log.Printf("⚠️ Error updating user_id %d: %v", id, err)
					continue
				}
				totalUpdated++
			}

			offset += batchSize
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"message":"All user embeddings updated successfully","count":%d}`, totalUpdated)))
		return nil
	}

	// --- Otherwise, update a single user's embedding (no password required) ---
	if userIDParam == "" {
		http.Error(w, "Missing userID", http.StatusBadRequest)
		return fmt.Errorf("missing userID")
	}

	var userID int64
	_, err := fmt.Sscan(userIDParam, &userID)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return fmt.Errorf("invalid userID: %v", err)
	}

	log.Printf("🔹 Updating embeddings for single user_id: %d", userID)

	err = vmdRepo.UpdateUserEmbeddingsFromVideoHistory(userID, 20)
	if err != nil {
		log.Printf("❌ Failed to update embeddings for user_id %d: %v", userID, err)
		http.Error(w, "Failed to update user embeddings", http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message":"Embeddings updated for user_id %d"}`, userID)))
	return nil
}


func main(){

	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	VideoMetaData := repository.NewPostgresVMDRepo(db)
	UserRepo := repository.NewPostgresUserRepo(db)
	

	log.Println("Server Started at Port 7110")
	http.HandleFunc("/update-user-embeddings", func(w http.ResponseWriter, r *http.Request) {
    	err := updateUserEmbeddings(w, r, UserRepo, VideoMetaData)
		if err != nil{
			log.Print(err)
		}
	})

	err := http.ListenAndServe(":7110", nil)
	if err != nil{
		log.Println("Critical Error Occured", "error", err)
		os.Exit(1)
	}
}