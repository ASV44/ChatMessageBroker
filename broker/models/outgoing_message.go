package models

import "time"

type OutgoingMessage struct {
	Channel string    `json:"channel"`
	Sender  string    `json:"user"`
	Text    string    `json:"text"`
	Time    time.Time `json:"time"`
}
