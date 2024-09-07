package models

type User struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	UserId   uint   `json:"user_id" gorm:"primaryKey;autoIncrement"`
	UserType int    `json:"user_type"`
}
