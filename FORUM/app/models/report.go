package models

type Report struct {
	User_id uint   `json:"user_id"`
	Post_id uint   `json:"post_id" gorm:"primaryKey;autoIncrement"`
	Content string `json:"content"`
	Reason  string `json:"reason"`
	Status  int    `json:"status"`
}
