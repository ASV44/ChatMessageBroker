package models

// User represents model with user data received from client in JSON format
type User struct {
	ID       int    `json:"ID"`
	NickName string `json:"nickName"`
}
