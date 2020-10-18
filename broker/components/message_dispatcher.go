package broker

import (
	"github.com/ASV44/ChatMessageBroker/broker/models"
)

type Dispatcher struct {
	workspace *Workspace
	ConnectionManager
	CommandDispatcher
	Transmitter
}

func NewMessageDispatcher(
	workspace *Workspace,
	connectionManager ConnectionManager,
	cmdDispatcher CommandDispatcher,
	transmitter Transmitter,
) Dispatcher {
	return Dispatcher{
		workspace:         workspace,
		ConnectionManager: connectionManager,
		CommandDispatcher: cmdDispatcher,
		Transmitter:       transmitter,
	}
}

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
			dispatcher.SendMessageToUser(sender, models.OutgoingMessage{Channel: channel.Name, Text: "not joined yet!"})
		}
	}
}
