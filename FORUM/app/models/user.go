package models

type User struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	User_id   uint   `json:"user_id" gorm:"primaryKey;autoIncrement"`
	User_type int    `json:"user_type"`
}
