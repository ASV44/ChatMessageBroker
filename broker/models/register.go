package models

import "time"

// Register represents model of message which is sent to user at connection for register
type Register struct {
	UserID int       `json:"userId"`
	Text   string    `json:"text"`
	Time   time.Time `json:"time"`
}
