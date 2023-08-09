package models

type Post struct{
	Id int32
	ParentId int32
	Content string
	CreateTime int32
	UserId int32
}