package models

import "time"

// OutgoingMessage represents model of message sent to client
type OutgoingMessage struct {
	Channel string    `json:"channel"`
	Sender  string    `json:"user"`
	Text    string    `json:"text"`
	Time    time.Time `json:"time"`
}
