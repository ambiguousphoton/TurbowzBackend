package main

import (
	"GoServer/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	_ "github.com/lib/pq"
)

func getEventMD(w http.ResponseWriter, r *http.Request, EventRepo repository.EventRepo) error {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return fmt.Errorf("wrong method: %s", r.Method)
	}

	eventIDStr := r.URL.Query().Get("event_id")
	if eventIDStr == "" {
		http.Error(w, "missing event_id", http.StatusBadRequest)
		return fmt.Errorf("missing event_id")
	}

	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid event_id", http.StatusBadRequest)
		return fmt.Errorf("invalid event_id: %v", err)
	}

	result, err := EventRepo.GetEventByID(eventID)
	if err != nil {
		http.Error(w, "error fetching event metadata: "+err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("error fetching event metadata: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return fmt.Errorf("failed to encode response: %v", err)
	}

	log.Println("Sent event metadata for event_id:", eventID)
	return nil
}

func getIndexPageEventIDs(w http.ResponseWriter, r *http.Request, EventRepo repository.EventRepo) error {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return fmt.Errorf("wrong method: %s", r.Method)
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

	events, err := EventRepo.GetTrendingEvents(limit, offset)
	if err != nil {
		http.Error(w, "error fetching trending events: "+err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("error fetching trending events: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return fmt.Errorf("failed to encode response: %v", err)
	}

	log.Printf("Sent %d trending events (offset: %d)\n", len(events), offset)
	return nil
}

func main() {
	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)

	EventRepo := repository.NewPostgresEventRepo(db)

	http.HandleFunc("/event-md", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if err := getEventMD(w, r, EventRepo); err != nil {
			log.Printf("Error in getEventMD: %v", err)
		}
	})

	http.HandleFunc("/index-events", func(w http.ResponseWriter, r *http.Request){
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if err := getIndexPageEventIDs(w, r, EventRepo); err != nil {
			log.Printf("Error in getIndexPageEventIDs: %v", err)
		}
	})

	log.Println("EventData Server Started at Port 7002")
	err := http.ListenAndServe(":7002", nil)
	if err != nil {
		log.Println("Critical Error Occured", "error", err)
		os.Exit(1)
	}
}