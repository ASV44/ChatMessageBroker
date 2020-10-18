package entity

import (
	"github.com/ASV44/ChatMessageBroker/common"
)

// User represents entity of broker user
type User struct {
	ID         int
	NickName   string
	Connection common.Connection
}
