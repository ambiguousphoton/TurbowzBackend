// Narayan Narayan Narayan Narayan÷
package models

import "database/sql"

type EcoMetaData struct {
	Eco_Id  			    int64
	Eco_Url 				string
	Uploader_ID 		    int64
	Uploader_Handle         string
	Uploader_Name           string
	Eco_Text       			string
	View_Count       		int64
	Luv_Count      			int64
	Comment_Count			int64
	Tags 					[]string
	Already_Luved			bool
	Images_Count			int
	Saves_Count				int64
	Created_At 				string
}

type EchoScore struct {
	Echo_Quality  sql.NullFloat64
	Echo_AI_Usage sql.NullFloat64
	Total_Qualtiy_Votes int64
	Total_Ai_Votes      int64
}