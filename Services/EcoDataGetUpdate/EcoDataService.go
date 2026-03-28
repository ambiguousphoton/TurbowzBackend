package main

import (
	"GoServer/repository"

	"log"
	"net/http"
	"os"
	"encoding/json"
	"strconv"
	_ "github.com/lib/pq" // postgres driver
	"GoServer/authenticator"
    "fmt"
)


func getEcoMD(w http.ResponseWriter, r *http.Request, EcoRepo repository.EcoRepo) error {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return fmt.Errorf("wrong method: %s", r.Method)
	}



	// --- Correct Query Param ---
	ecoIDStr := r.URL.Query().Get("eco_id")
	if ecoIDStr == "" {
		http.Error(w, "missing eco_id", http.StatusBadRequest)
		return fmt.Errorf("missing eco_id")
	}

	ecoID, err := strconv.ParseInt(ecoIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid eco_id", http.StatusBadRequest)
		return fmt.Errorf("invalid eco_id: %v", err)
	}

	// --- Fetch Metadata ---
	result, err := EcoRepo.GetEcoMetaData(ecoID)
	if err != nil {
		http.Error(w, "error fetching eco metadata: "+err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("error fetching eco metadata: %v", err)
	}

	// JSON Output
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return fmt.Errorf("failed to encode response: %v", err)
	}

	log.Println("Sent eco metadata for eco_id:", ecoID)
	return nil
}



// func ecoViewUpdate(w http.ResponseWriter, r *http.Request, EcoRepo repository.EcoRepo, userRepo repository.UserRepo) {
//     log.Println("videoViewUpdate called") 

//     w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
//     w.Header().Set("Access-Control-Allow-Methods", "POST")
//     w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

//     if r.Method != http.MethodPost {
//         log.Printf("Wrong Method: %s\n", r.Method)
//         http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//         return
//     }

//     if err := r.ParseForm(); err != nil {
//         log.Printf("Failed to parse form: %v\n", err)
//         http.Error(w, "failed to parse form", http.StatusBadRequest)
//         return
//     }

//     videoIDStr := r.FormValue("video_id")
//     userIDStr  := r.FormValue("user_id")
//     log.Printf("Received video_id: %s\n", videoIDStr)

//     videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
//     if err != nil {
//         log.Printf("Invalid video_id: %s, error: %v\n", videoIDStr, err)
//         http.Error(w, "invalid video_id", http.StatusBadRequest)
//         return
//     }
//     if userIDStr == "" {
//         log.Printf("Not able to record History, User not signed In, no userID")
//     }else {
//         userID, err := strconv.ParseInt(userIDStr, 10, 64)
//         if err != nil {
//             log.Printf("Invalid userID: %s, error: %v\n", userIDStr, err)
//             http.Error(w, "invalid userID", http.StatusBadRequest)
//         } else{
//             err = userRepo.AddVideoInUserHistory(userID, videoID)
//             if err != nil{
//                 log.Printf("Error Updating UserHistory %v", err)
//             }
//             log.Printf("Video add in User History")
//         }
//     }

//     err = vmdRepo.VideoViewUpdate(videoID)
//     if err != nil {
//         log.Printf("Error updating views for video_id %d: %v\n", videoID, err)
//         http.Error(w, "failed to update views", http.StatusInternalServerError)
//         return
//     }

//     log.Printf("Views updated successfully for video_id %d\n", videoID)
//     w.WriteHeader(http.StatusOK)
//     w.Write([]byte("views updated"))
// }


func updateLuv(w http.ResponseWriter, r *http.Request, EcoRepo repository.EcoRepo) {
    log.Println("updateLuv for eco called") 

    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method != http.MethodPost {
        log.Printf("Wrong Method: %s\n", r.Method)
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    if err := r.ParseForm(); err != nil {
        log.Printf("Failed to parse form: %v\n", err)
        http.Error(w, "failed to parse form", http.StatusBadRequest)
        return
    }
	userID, ok := r.Context().Value("userID").(int64)
	if !ok  {
		log.Printf("updateLuv: Invalid or missing userID in context")
		http.Error(w, "Error InvalidUserId ", http.StatusBadRequest)
		return 
	}


    ecoIDStr := r.FormValue("eco_id")
    log.Printf("Received eco_id: %s\n", ecoIDStr)

    ecoID, err := strconv.ParseInt(ecoIDStr, 10, 64)
    if err != nil {
        log.Printf("Invalid eco_id: %s, error: %v\n", ecoIDStr, err)
        http.Error(w, "invalid eco_id", http.StatusBadRequest)
        return
    }

    luved, err := EcoRepo.UpdateLuv(ecoID, userID)
    if err != nil {
        log.Printf("Error updating luvs for eco_id %d and user_id %d: %v\n", ecoID, userID, err)
        http.Error(w, "failed to update luv", http.StatusInternalServerError)
        return
    }

    log.Printf("luv updated successfully for eco_id %d\n", ecoID)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]bool{"luved": luved})
}

func ecoLuvStatus(w http.ResponseWriter, r *http.Request, EcoRepo repository.EcoRepo) {
    log.Println("ecoLuvStatus called")

    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method != http.MethodPost {
        log.Printf("Wrong Method: %s\n", r.Method)
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    if err := r.ParseForm(); err != nil {
        log.Printf("Failed to parse form: %v\n", err)
        http.Error(w, "failed to parse form", http.StatusBadRequest)
        return
    }

    userIDStr := r.FormValue("user_ID")
    userID, err := strconv.ParseInt(userIDStr, 10, 64)
    if err != nil {
        log.Printf("ecoLuvStatus: Invalid or missing userID")
        http.Error(w, "invalid user_ID", http.StatusBadRequest)
        return
    }

    ecoIDStr := r.FormValue("eco_id")
    ecoID, err := strconv.ParseInt(ecoIDStr, 10, 64)
    if err != nil {
        log.Printf("Invalid eco_id: %s, error: %v\n", ecoIDStr, err)
        http.Error(w, "invalid eco_id", http.StatusBadRequest)
        return
    }

    // 🔍 Get luv status + total luvs
    luved, totalLuvs, err := EcoRepo.LuvStatus(ecoID, userID)
    if err != nil {
        log.Printf("Error fetching luv status for eco_id %d, user_id %d: %v\n", ecoID, userID, err)
        http.Error(w, "failed to fetch luv status", http.StatusInternalServerError)
        return
    }

    log.Printf("luv status fetched successfully for eco_id %d\n", ecoID)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "luved":      luved,
        "total_luvs": totalLuvs,
    })
}

func getTrendingEcos(w http.ResponseWriter, r *http.Request, EcoRepo repository.EcoRepo) {
w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	videos, err := EcoRepo.GetTrendingEcos(limit, offset)
	if err != nil {
		log.Printf("Error fetching trending ecos %d: %v", err)
		http.Error(w, "Failed to get trending ecos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

func getEchoScore(w http.ResponseWriter, r *http.Request, echoRepo repository.EcoRepo){
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("Wrong Method")
		return 
	}
    if err := r.ParseForm(); err != nil {
        log.Printf("Failed to parse form: %v\n", err)
        http.Error(w, "failed to parse form", http.StatusBadRequest)
        return  
    }

	echo_id_str := r.URL.Query().Get("echo_id")
	echo_id, err := strconv.ParseInt(echo_id_str, 10, 64)
	if err != nil{
        http.Error(w, "invalid echo_id", http.StatusBadRequest)
		return
	}

	echo_score, err := echoRepo.GetEchoScore(echo_id)
	if err != nil {
        log.Printf("Error Getting Score for echo_id %d: %v\n", echo_id, err)
        http.Error(w, "failed to get score", http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(echo_score)	
	if err != nil{
		log.Printf("Error Sending the Echo Score for Echo ID %d, getting error %v", echo_id, err)
	}
	log.Println("Successfuly sent the Echo Score for Echo ID %d", echo_id)
}

func main() {
	
	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	EcoRepo := repository.NewPostgresEcoRepo(db)
    // UserRepo :=  repository.NewPostgresUserRepo(db)

	http.HandleFunc("/emd", func(w http.ResponseWriter, r *http.Request) {
    	err := getEcoMD(w, r, EcoRepo)
		if err != nil {
			http.Error(w, "Failed to get Eco Meta Data", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/luv", authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
    	updateLuv(w, r, EcoRepo)
	}))

	// http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
    // 	ecoViewUpdate(w, r, EcoRepo, UserRepo)	
	// })

    http.HandleFunc("/check-eco-luv-status", func(w http.ResponseWriter, r *http.Request){
        ecoLuvStatus(w, r, EcoRepo)
    })

    http.HandleFunc("/get-trending-ecos", func(w http.ResponseWriter, r *http.Request) {
        getTrendingEcos(w, r, EcoRepo)
    })

	http.HandleFunc("/get-echo-score", func(w http.ResponseWriter, r *http.Request) {
        getEchoScore(w, r, EcoRepo)
    })

	log.Println("EcoData GETER Server Started at Port 7011")
	err := http.ListenAndServe(":7011", nil)
	if err != nil{
		log.Println("Critical Error Occured", "error", err)
		os.Exit(1)
	}
}
