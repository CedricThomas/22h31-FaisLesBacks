package model

import "time"

type CreateReminderRequest struct {
	Title   string    `json:"title" binding:"required"`
	Content string    `json:"content"`
	Date    time.Time `json:"date" binding:"required"`
}

type Reminder struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Date      time.Time `json:"date"`
	Triggered bool      `json:"triggered"`
}
