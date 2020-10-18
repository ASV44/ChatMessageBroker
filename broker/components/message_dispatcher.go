package broker

import (
	"github.com/ASV44/ChatMessageBroker/broker/entity"
	"github.com/ASV44/ChatMessageBroker/broker/models"
	"github.com/ASV44/ChatMessageBroker/broker/services"
)

// Dispatcher represents broker component which process new incoming message from client
type Dispatcher struct {
	workspace *Workspace
	ConnectionManager
	CommandDispatcher
	services.Transmitter
}

// NewDispatcher creates new instance of Dispatcher
func NewDispatcher(
	workspace *Workspace,
	connectionManager ConnectionManager,
	cmdDispatcher CommandDispatcher,
	transmitter services.Transmitter,
) Dispatcher {
	return Dispatcher{
		workspace:         workspace,
		ConnectionManager: connectionManager,
		CommandDispatcher: cmdDispatcher,
		Transmitter:       transmitter,
	}
}

// DispatchMessage process incoming message by type and invoke specific method
func (dispatcher Dispatcher) DispatchMessage(message models.IncomingMessage) {
	switch message.Type {
	case models.CMD:
		dispatcher.DispatchCommand(message)
	case models.DIRECT:
		dispatcher.handleDirectMessage(message)
	case models.CHANNEL:
		dispatcher.handleChannelMessage(message)
	}
}

func (dispatcher Dispatcher) handleDirectMessage(message models.IncomingMessage) {
	user, isPresent := dispatcher.workspace.users[message.Target]
	if isPresent {
		dispatcher.SendMessageToUser(user, message.ToOutgoingMessage())
	}
}

func (dispatcher Dispatcher) handleChannelMessage(message models.IncomingMessage) {
	if channel, isPresent := dispatcher.workspace.channels[message.Target]; isPresent {
		sender := dispatcher.workspace.users[message.Sender.NickName]
		if channel.Contains(sender) {
			draft := message.ToOutgoingMessage()
			for _, user := range channel.Subscribers {
				if user.ID != message.Sender.ID {
					dispatcher.SendMessageToUser(user, draft)
				}
			}
		} else {
			dispatcher.SendMessageToUser(
				sender,
				models.OutgoingMessage{
					Channel: channel.Name,
					Text:    entity.ChannelNotJoined{Name: channel.Name}.Error(),
				},
			)
		}
	}
}
