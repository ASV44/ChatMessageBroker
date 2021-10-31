package entity

import "github.com/ASV44/chat-message-broker/common"

// RegistrationData represents entity which combine all account data for registration with connection
type RegistrationData struct {
	NickName     string
	PasswordHash string
	Connection   common.Connection
}
