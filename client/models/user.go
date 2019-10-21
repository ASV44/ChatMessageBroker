package models

// User represents model with user data which will be sent to broker
type User struct {
	ID       int    `json:"id"`
	NickName string `json:"nickName"`
}
