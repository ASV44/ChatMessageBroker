package receiver

import "time"

// Register represents model of message received from broker at registering user at broker
type Register struct {
	Text string    `json:"text"`
	Time time.Time `json:"time"`
}

// AccountData represents model of message which is sent from user with all account data required at sign up
type AccountData struct {
	NickName string `json:"nickName"`
}
