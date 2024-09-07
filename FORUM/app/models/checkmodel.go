package models

type Check struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	Reason   string `json:"reason"`
	PostId  uint   `json:"post_id"`
}
