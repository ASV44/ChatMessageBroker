package broker

import (
	"fmt"
	"github.com/ASV44/ChatMessageBroker/broker/entity"
	"github.com/ASV44/ChatMessageBroker/broker/models"
)

type Transmitter interface {
	SendMessageToUser(entity.User, models.OutgoingMessage)
}

type CommunicationManager struct{}

func NewCommunicationManager() CommunicationManager {
	return CommunicationManager{}
}

func (manager CommunicationManager) SendMessageToUser(user entity.User, message models.OutgoingMessage) {
	err := user.Connection.SendMessage(message)
	if err != nil {
		fmt.Printf("Failed to send message to user ID:%s Nickname: %s. Error: %s\n", user.ID, user.NickName, err)
		// Handle error of sending message here.
		// Here is possible to implement various scenarios of handling communication error.
		// Resending of message, saving in DB, adding in queue of not delivered messages.
	}
}
