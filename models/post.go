package models

type Post struct{
	ID int `json:"id"`
	ParentID int `json:"parent_id"`
	Content string `json:"content"`
	CreateTime int `json:"create_time"`
	UserID int `json:"user_id"`
}

