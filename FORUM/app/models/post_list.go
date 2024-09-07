package models

import "time"

type PostList struct {
	Content string    `json:"content"`
	UserId uint      `json:"use_id"`
	Id      uint      `json:"id"`
	Time    time.Time `json:"time"`
}
