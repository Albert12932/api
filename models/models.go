package models

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

type Task struct {
	Id     int    `json:"id"`
	UserID int    `json:"userId" binding:"required"`
	Header string `json:"header" binding:"required"`
	Text   string `json:"text" binding:"required"`
	Done   bool   `json:"done" default:"false"`
	// TODO date_created
}

//TODO API response
