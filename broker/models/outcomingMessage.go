package models

import "time"

type OutcomingMessage struct {
	Channel string    `json:"channel"`
	Sender  string    `json:"user"`
	Text    string    `json:"text"`
	Time    time.Time `json:"time"`
}
