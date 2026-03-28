package models

import "database/sql"
// import     "github.com/lib/pq"

type VideoMetaData struct {
	Video_ID  			  string
	Uploader_ID 		   int64
	Uploader_Handle       string
	Uploader_Name         string
	Uploader_Url          string
	Title       			string
	Video_Info 				string
	Video_Url   			string   /// Internal URL
	Views       			int64
	Luvs      			    int64
	Upload_Time 			string
	Tags 					[]string
	Already_Luved			  bool
	Watched_At 				string
}


type VideoScore struct {
	Video_Quality  sql.NullFloat64
	Video_AI_Usage sql.NullFloat64
	Total_Qualtiy_Votes int64
	Total_Ai_Votes      int64
}