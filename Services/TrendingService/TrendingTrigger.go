package main

import (
	"log"
	"math"
	"time"

	"GoServer/repository"
)

func init() {
	// Better logs
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// trending score calculation
func trendingScore(views, likes, comments, shares int64, hoursSince float64) float64 {
	currentScore := (0.5 * float64(views)) +
		(1.0 * float64(likes)) +
		(1.5 * float64(comments)) +
		(2.0 * float64(shares))

	k := 0.03 // decay constant

	return currentScore * math.Exp(-k*hoursSince)
}

func parseTimestamp(timestampStr string) (time.Time, error) {

	// If DB returned an empty string or NULL
	if timestampStr == "" {
		return time.Time{}, nil
	}

	// Try RFC3339
	t, err := time.Parse(time.RFC3339, timestampStr)
	if err == nil {
		return t, nil
	}

	// PostgreSQL default format
	formats := []string{
		"2006-01-02 15:04:05",               // without timezone
		"2006-01-02 15:04:05.999999",       // with microseconds
		"2006-01-02 15:04:05-07",           // with timezone
		"2006-01-02 15:04:05.999999-07",    // microseconds + tz
	}

	for _, f := range formats {
		t, err2 := time.Parse(f, timestampStr)
		if err2 == nil {
			return t, nil
		}
	}

	return time.Time{}, err
}

func PushTrendingScore(videoID int64, vmdRepo repository.VMDrepo) error {

	views, likes, comments, shares, lastScore, timestampStr, err :=
		vmdRepo.GetEngagementInfo(videoID)

	if err != nil {
		log.Println("GetEngagementInfo error:", err)
		return err
	}

	parsedTime, err := parseTimestamp(timestampStr)
	if err != nil {
		log.Println("Timestamp parse error:", timestampStr, err)
		return err
	}

	// If timestamp is invalid or zero → decayed fully

	hoursSince := time.Since(parsedTime).Hours()


	currentScore := trendingScore(views, likes, comments, shares, hoursSince)
	delta := currentScore - lastScore

	err = vmdRepo.SaveTrendingDelta(videoID, currentScore, delta)
	if err != nil {
		log.Println("SaveTrendingDelta error:", err)
		return err
	}

	return nil
}

func StartTrendingVideoScheduler(vmdRepo repository.VMDrepo) {

	go func() {
		for {
			log.Println("Running trending score update...")

			limit := 500
			offset := 0

			for {
				ids, err := vmdRepo.GetVideoIDsPaginated(limit, offset)
				if err != nil {
					log.Println("Error fetching video IDs:", err)
					break
				}

				if len(ids) == 0 {
					break
				}

				log.Println("Processing batch, count:", len(ids), "offset:", offset)

				for _, id := range ids {
					err := PushTrendingScore(id, vmdRepo)
					if err != nil {
						log.Println("Error updating trending for video:", id, err)
					}
				}

				offset += limit
			}

			log.Println("Trending update finished. Sleeping for 1 hour.")
			time.Sleep(1 * time.Hour)
		}
	}()
}


func PushEcoTrendingScore(ecoID int64, ecoRepo repository.EcoRepo) error {

    views, likes, comments, shares, lastScore, timestampStr, err :=
        ecoRepo.GetEcoEngagementInfo(ecoID)

    if err != nil {
        log.Println("Eco GetEngagementInfo error:", err)
        return err
    }

    parsedTime, err := parseTimestamp(timestampStr)
    if err != nil {
        log.Println("Eco timestamp parse error:", timestampStr, err)
        return err
    }

    hoursSince := time.Since(parsedTime).Hours()

    currentScore := trendingScore(views, likes, comments, shares, hoursSince)
    delta := currentScore - lastScore

    err = ecoRepo.SaveEcoTrendingDelta(ecoID, currentScore, delta)
    if err != nil {
        log.Println("Eco SaveTrendingDelta error:", err)
        return err
    }

    return nil
}

func StartEcoTrendingScheduler(ecoRepo repository.EcoRepo) {

    go func() {
        for {
            log.Println("Running ECO trending score update...")

            limit := 500
            offset := 0

            for {
                ids, err := ecoRepo.GetAllEcoIDs(limit, offset)
                if err != nil {
                    log.Println("Eco error fetching Eco IDs:", err)
                    break
                }

                if len(ids) == 0 {
                    break
                }

                log.Println("Eco Processing batch:", len(ids), "offset:", offset)

                for _, id := range ids {
                    err := PushEcoTrendingScore(id, ecoRepo)
                    if err != nil {
                        log.Println("Eco error updating trending:", id, err)
                    }
                }

                offset += limit
            }

            log.Println("ECO trending update finished. Sleeping for 1 hour.")
            time.Sleep(1 * time.Hour)
        }
    }()
}

func PushEventTrendingScore(eventID int64, eventRepo repository.EventRepo) error {

    views, likes, comments, shares, lastScore, timestampStr, err :=
        eventRepo.GetEventEngagementInfo(eventID)

    if err != nil {
        log.Println("Event GetEngagementInfo error:", err)
        return err
    }

    parsedTime, err := parseTimestamp(timestampStr)
    if err != nil {
        log.Println("Event timestamp parse error:", timestampStr, err)
        return err
    }

    hoursSince := time.Since(parsedTime).Hours()

    currentScore := trendingScore(views, likes, comments, shares, hoursSince)
    delta := currentScore - lastScore

    err = eventRepo.SaveEventTrendingDelta(eventID, currentScore, delta)
    if err != nil {
        log.Println("Event SaveTrendingDelta error:", err)
        return err
    }

    return nil
}

func StartEventTrendingScheduler(eventRepo repository.EventRepo) {

    go func() {
        for {
            log.Println("Running Event trending score update...")

            limit := 500
            offset := 0

            for {
                ids, err := eventRepo.GetAllEventIDs(limit, offset)
                if err != nil {
                    log.Println("Event error fetching Event IDs:", err)
                    break
                }

                if len(ids) == 0 {
                    break
                }

                log.Println("Event Processing batch:", len(ids), "offset:", offset)

                for _, id := range ids {
                    err := PushEventTrendingScore(id, eventRepo)
                    if err != nil {
                        log.Println("Event error updating trending:", id, err)
                    }
                }

                offset += limit
            }

            log.Println("Event trending update finished. Sleeping for 1 hour.")
            time.Sleep(1 * time.Hour)
        }
    }()
}


func main() {

	log.Println("Starting Triggers...")

	dbDestination := "host=localhost port=5454 user=postgres password=Narayan!123 dbname=MetaDataStorage sslmode=disable"
	db := repository.NewPostgresDB(dbDestination)
	vmdRepo := repository.NewPostgresVMDRepo(db)
	ecoRepo := repository.NewPostgresEcoRepo(db)
	eventRepo := repository.NewPostgresEventRepo(db)

	log.Println("Trending Scheduler Started")
	StartTrendingVideoScheduler(vmdRepo)
	StartEcoTrendingScheduler(ecoRepo)
	StartEventTrendingScheduler(eventRepo)

	select {} // keep app running
}
