package models

import "time"

type Register struct {
	UserId int       `json:"userId"`
	Text   string    `json:"text"`
	Time   time.Time `json:"time"`
}
