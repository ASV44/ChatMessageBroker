package entity

import "github.com/ASV44/ChatMessageBroker/common"

// RegistrationData represents entity which combine all account data for registration with connection
type RegistrationData struct {
	NickName   string
	Connection common.Connection
}
