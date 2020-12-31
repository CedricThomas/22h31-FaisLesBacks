package model

import "time"

type Location struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

type CreateMemoRequest struct {
	Title    string    `json:"title" binding:"required"`
	Content  string    `json:"content" binding:"required"`
	Location *Location `json:"location"`
}

type UpdateMemoRequest struct {
	CreateMemoRequest
}

type Memo struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Location  string    `json:"location"`
}
