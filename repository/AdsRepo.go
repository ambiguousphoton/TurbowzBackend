package repository

import (
	"GoServer/models"
	"database/sql"
)

type AdsRepo interface {
	CreateNewBannerAd(ad *models.BannerAd) (int64, error)
	GetBannerAds(page, limit int) ([]models.BannerAd, error)
}

type PostgresAdsRepo struct {
	db *sql.DB
}

func NewPostgresAdsRepo(db *sql.DB) AdsRepo {
	return &PostgresAdsRepo{db: db}
}

func (r *PostgresAdsRepo) CreateNewBannerAd(ad *models.BannerAd) (int64, error) {
	query := `
		INSERT INTO banner_ads (uploader_id, title, redirect_url, start_date, end_date, views, clicks, created_at)
		VALUES ($1, $2, $3, NOW(), NOW() + INTERVAL '30 days', 0, 0, NOW())
		RETURNING ad_id;
	`

	var insertedID int64
	err := r.db.QueryRow(query, ad.UploaderID, ad.Title, ad.RedirectURL).Scan(&insertedID)
	if err != nil {
		return 0, err
	}

	return insertedID, nil
}

func (r *PostgresAdsRepo) GetBannerAds(page, limit int) ([]models.BannerAd, error) {
	offset := (page - 1) * limit

	query := `
		SELECT ad_id, uploader_id, title, redirect_url, created_at
		FROM banner_ads
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ads := []models.BannerAd{}

	for rows.Next() {
		var ad models.BannerAd
		err := rows.Scan(&ad.AdID, &ad.UploaderID, &ad.Title, &ad.RedirectURL, &ad.CreatedAt)
		if err != nil {
			return nil, err
		}

		ads = append(ads, ad)
	}

	return ads, nil
}
