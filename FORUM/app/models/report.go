package models

type Report struct {
	UserId uint   `json:"user_id"`
	PostId uint   `json:"post_id" gorm:"primaryKey;autoIncrement"`
	Reason  string `json:"reason"`
	Status  int    `json:"status"`
}
