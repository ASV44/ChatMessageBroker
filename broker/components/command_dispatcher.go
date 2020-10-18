package broker

import (
	"github.com/ASV44/ChatMessageBroker/broker/entity"
	"github.com/ASV44/ChatMessageBroker/broker/models"
	"github.com/ASV44/ChatMessageBroker/broker/services"
	"strings"
	"time"
)

// Available commands
const (
	CreateChannel = "create"
	JoinChannel   = "join"
	LeaveChannel  = "leave"
	Show          = "show"
)

// Show command options
const (
	Users    = "users"
	Channels = "channels"
	All      = "all"
)

type CommandDispatcher struct {
	workspace *Workspace
	services.Transmitter
}

func NewCommandDispatcher(workspace *Workspace, transmitter services.Transmitter) CommandDispatcher {
	return CommandDispatcher{workspace: workspace, Transmitter: transmitter}
}

func (dispatcher CommandDispatcher) DispatchCommand(message models.IncomingMessage) {
	switch message.Target {
	case CreateChannel:
		dispatcher.createChannel(message.Sender, message.Text)
	case JoinChannel:
		dispatcher.joinChannel(message.Sender, message.Text)
	case LeaveChannel:
		dispatcher.leaveChannel(message.Sender, message.Text)
	case Show:
		dispatcher.show(message.Sender, message.Text)
	}
}

func (dispatcher CommandDispatcher) createChannel(sender models.User, name string) {
	user := dispatcher.workspace.users[sender.NickName]
	err := dispatcher.workspace.CreateChannel(sender, name)
	if err != nil {
		dispatcher.SendMessageToUser(user, models.OutgoingMessage{Channel: name, Text: err.Error()})
		return
	}

	dispatcher.SendMessageToUser(user, dispatcher.getWorkspaceChannelsMessage())
}

func (dispatcher CommandDispatcher) joinChannel(sender models.User, name string) {
	user := dispatcher.workspace.users[sender.NickName]
	err := dispatcher.workspace.AddUserToChannel(sender, name)
	if err != nil {
		dispatcher.SendMessageToUser(user, models.OutgoingMessage{Channel: name, Text: err.Error()})
	}

	dispatcher.SendMessageToUser(user, dispatcher.getChannelSubscribersMessage(dispatcher.workspace.channels[name]))
}

func (dispatcher CommandDispatcher) leaveChannel(sender models.User, name string) {
	user := dispatcher.workspace.users[sender.NickName]
	err := dispatcher.workspace.RemoveUserFromChannel(sender, name)
	if err != nil {
		dispatcher.SendMessageToUser(user, models.OutgoingMessage{Channel: name, Text: err.Error()})
	}

	dispatcher.SendMessageToUser(user, dispatcher.getChannelSubscribersMessage(dispatcher.workspace.channels[name]))
}

func (dispatcher CommandDispatcher) show(sender models.User, param string) {
	user := dispatcher.workspace.users[sender.NickName]
	switch param {
	case Users:
		dispatcher.SendMessageToUser(user, dispatcher.getWorkspaceUsersMessage())
	case Channels:
		dispatcher.SendMessageToUser(user, dispatcher.getWorkspaceChannelsMessage())
	case All:
		dispatcher.SendMessageToUser(user, dispatcher.getWorkspaceChannelsMessage())
		dispatcher.SendMessageToUser(user, dispatcher.getWorkspaceUsersMessage())
	default: // Get all users of specific channel if it exist
		if channel, exist := dispatcher.workspace.channels[param]; exist {
			dispatcher.SendMessageToUser(user, dispatcher.getChannelSubscribersMessage(channel))
		}
	}
}

func (dispatcher CommandDispatcher) getWorkspaceChannelsMessage() models.OutgoingMessage {
	return models.OutgoingMessage{
		Channel: "Channels",
		Text:    strings.Join(dispatcher.workspace.WorkspaceChannels(), " | "),
		Time:    time.Now(),
	}
}

func (dispatcher CommandDispatcher) getWorkspaceUsersMessage() models.OutgoingMessage {
	return models.OutgoingMessage{
		Sender: "Users",
		Text:   strings.Join(dispatcher.workspace.WorkspaceUsers(), " | "),
		Time:   time.Now(),
	}
}

func (dispatcher CommandDispatcher) getChannelSubscribersMessage(channel entity.Channel) models.OutgoingMessage {
	return models.OutgoingMessage{
		Channel: channel.Name,
		Sender:  "Users",
		Text:    strings.Join(dispatcher.workspace.ChannelSubscribers(channel), " | "),
		Time:    time.Now(),
	}
}
