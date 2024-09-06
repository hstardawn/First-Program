package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	Content string    `json:"content"`
	User_id uint      `json:"user_id"`
	Post_id uint      `json:"post_id" gorm:"primaryKey;autoIncrement"`
	Time    time.Time `json:"time"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
