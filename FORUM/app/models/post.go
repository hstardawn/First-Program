package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	Content string    `json:"content"`
	UserId uint      `json:"user_id"`
	PostId uint      `json:"post_id" gorm:"primaryKey;autoIncrement"`
	Time    time.Time `json:"time"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
