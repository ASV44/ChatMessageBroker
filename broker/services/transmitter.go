package services

import (
	"fmt"
	"github.com/ASV44/chat-message-broker/broker/entity"
	"github.com/ASV44/chat-message-broker/broker/models"
)

// Transmitter represents abstraction for performing sending of message to client
type Transmitter interface {
	SendMessageToUser(entity.User, models.OutgoingMessage)
}

// CommunicationManager represents service for performing sending of message to client
type CommunicationManager struct{}

// NewCommunicationManager creates new instance of CommunicationManager
func NewCommunicationManager() CommunicationManager {
	return CommunicationManager{}
}

// SendMessageToUser perform action of sending message to user and handle error
func (manager CommunicationManager) SendMessageToUser(user entity.User, message models.OutgoingMessage) {
	err := user.Connection.SendMessage(message)
	if err != nil {
		fmt.Printf("Failed to send message to user ID:%d Nickname: %s. Error: %s\n", user.ID, user.NickName, err)
		// Handle error of sending message here.
		// Here is possible to implement various scenarios of handling communication error.
		// Resending of message, saving in DB, adding in queue of not delivered messages.
	}
}
