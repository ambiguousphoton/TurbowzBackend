package models


type VideoSearchResult struct {
    VideoID    int
    UploaderName  string
    UploaderHandle string
    Title  string
    Views int
	VideoURL   string
    Date string
    Tags []string
}

// type VideoSearchResult struct {
//     VideoID    int    `json:"videoID"`
//     UploaderID string `json:"uploaderID"`
//     Title      string `json:"title"`
//     Views      int    `json:"views"`
//     VideoURL   string `json:"videoURL"`
// }