package components

import (
	"fmt"
	"github.com/ASV44/ChatMessageBroker/client/models/receiver"
	"github.com/ASV44/ChatMessageBroker/client/models/sender"
	"github.com/ASV44/ChatMessageBroker/common"
	"io"
)

// CommunicationService represents abstraction of for performing message communication with broker
type CommunicationService interface {
	SendMessage(sender.Message)
	GetMessage() (receiver.Message, error)
}

// CommunicationManager represents service for performing message communication with broker
type CommunicationManager struct {
	connection common.Connection
}

// NewCommunicationManager creates new instance of CommunicationManager
func NewCommunicationManager(connection common.Connection) CommunicationManager {
	return CommunicationManager{connection: connection}
}

// SendMessage send message to broker and handle error in case when sending of message fails
func (manager CommunicationManager) SendMessage(message sender.Message) {
	err := manager.connection.SendMessage(message)
	if err != nil {
		fmt.Printf(
			"Failed to send message.\nType:%s\nTarget: %s\nTime: %s\nErr: %s\n",
			message.Type,
			message.Target,
			message.Time,
			err,
		)
		// Handle error of sending message here.
		// Here is possible to implement various scenarios of handling communication error.
		// Resending of message, saving in DB, adding in queue of not delivered messages.
	}
}

// GetMessage get from broker connection and handle error
func (manager CommunicationManager) GetMessage() (receiver.Message, error) {
	var message receiver.Message
	if err := manager.connection.GetMessage(&message); err != nil && err != io.EOF {
		return receiver.Message{}, err
	}

	return message, nil
}
