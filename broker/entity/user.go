package entity

import (
	broker "github.com/ASV44/ChatMessageBroker/broker/components"
)

// User represents entity of broker user
type User struct {
	ID         int
	NickName   string
	Connection broker.Connection
}
