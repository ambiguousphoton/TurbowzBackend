package models

type VideoMetaData struct {
	Video_ID  			  string
	Uploader_ID 		   int64
	Uploader_Handle       string
	Uploader_Name         string
	Title       			string
	Video_Info 				string
	Video_Url   			string   /// Internal URL
	Views       			int64
	Agrees      			int64
	Disagrees   			int64
	Upload_Time 			string
}
