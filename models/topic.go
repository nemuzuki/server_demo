package models

type Topic struct{
	ID int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	CreateTime int `json:"create_time"`
}

