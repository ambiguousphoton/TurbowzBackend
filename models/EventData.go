// Narayan Narayan Narayan Narayan 

// Narayan Narayan Narayan Narayan÷
package models

type EventMetaData struct {
	Event_Id  			    int64
	Event_Url 				string
	Event_Title				string
	Uploader_ID 		    int64
	Uploader_Handle         string
	Uploader_Name           string
	Event_Description 		string
	View_Count       		int64
	Luv_Count      			int64
	Comment_Count			int64
	Tags 					[]string
	Already_Luved			bool
	Images_Count			int
	Saves_Count				int64
	Created_At 				string
	Event_Start_Time        string
	Event_End_Time          string
}
