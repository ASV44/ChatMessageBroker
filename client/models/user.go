package models

// User represents model with user data which will be sent to broker
type User struct {
	ID       int    `json:"ID"`
	NickName string `json:"nickName"`
}
