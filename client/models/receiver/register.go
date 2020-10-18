package receiver

import "time"

// Register represents model of message received from broker at registering user at broker
type Register struct {
	UserID int       `json:"userId"`
	Text   string    `json:"text"`
	Time   time.Time `json:"time"`
}
