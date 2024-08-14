package models

import "time"

type PostList struct {
	Content string    `json:"content"`
	User_id uint      `json:"user_id"`
	Id      uint      `json:"id"`
	Time    time.Time `json:"time"`
}
