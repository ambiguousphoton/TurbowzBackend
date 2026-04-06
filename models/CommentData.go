package models
import "database/sql"

type CommentData struct {
	Comment_id int64
	Commenter_id int64
	Parent_video_id int64
	Comment_text string
	Comment_date string
	Commenter_Handle string
	Commenter_Name string
	Parent_Comment_ID sql.NullInt64
}


type EcoCommentData struct {
	Comment_id int64
	Commenter_id int64
	Parent_Eco_id int64
	Comment_text string
	Comment_date string
	Commenter_Handle string
	Commenter_Name string
	Parent_Comment_ID sql.NullInt64
}