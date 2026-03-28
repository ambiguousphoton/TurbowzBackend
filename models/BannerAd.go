package models

type BannerAd struct {
	AdID        int64  `json:"ad_id"`
	UploaderID  int64  `json:"uploader_id"`
	Title       string `json:"title"`
	RedirectURL string `json:"redirect_url"`
	CreatedAt   string `json:"created_at"`
}
