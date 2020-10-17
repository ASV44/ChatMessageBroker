package broker

import (
	"github.com/ASV44/ChatMessageBroker/broker"
	"github.com/ASV44/ChatMessageBroker/broker/models"
)

type MessageDispatcher struct {
	workspace     *Workspace
	cmdDispatcher CommandDispatcher
}

func NewMessageDispatcher(workspace *Workspace) MessageDispatcher {
	return MessageDispatcher{
		workspace:     workspace,
		cmdDispatcher: NewCommandDispatcher(workspace),
	}
}

func (dispatcher MessageDispatcher) DispatchMessage(message models.IncomingMessage) {
	switch message.Type {
	case models.CMD:
		dispatcher.cmdDispatcher.DispatchCommand(message)
	case models.DIRECT:
		dispatcher.handleDirectMessage(message)
	case models.CHANNEL:
		dispatcher.handleChannelMessage(message)
	}
}

func (dispatcher MessageDispatcher) handleDirectMessage(message models.IncomingMessage) {
	user, isPresent := dispatcher.workspace.users[message.Target]
	if isPresent {
		broker.SendMessageToUser(user, message.ToOutgoingMessage())
	}
}

func (dispatcher MessageDispatcher) handleChannelMessage(message models.IncomingMessage) {
	if channel, isPresent := dispatcher.workspace.channels[message.Target]; isPresent {
		sender := dispatcher.workspace.users[message.Sender.NickName]
		if channel.Contains(sender) {
			draft := message.ToOutgoingMessage()
			for _, user := range channel.Subscribers {
				if user.ID != message.Sender.ID {
					broker.SendMessageToUser(user, draft)
				}
			}
		} else {
			broker.SendMessageToUser(sender, models.OutgoingMessage{Channel: channel.Name, Text: "not joined yet!"})
		}
	}
}
