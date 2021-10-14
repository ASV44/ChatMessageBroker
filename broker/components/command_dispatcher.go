package broker

import (
	"github.com/ASV44/chat-message-broker/broker/entity"
	"github.com/ASV44/chat-message-broker/broker/models"
	"github.com/ASV44/chat-message-broker/broker/services"
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

// CommandDispatcher represents broker component which process new incoming command from client
type CommandDispatcher struct {
	workspace *Workspace
	services.Transmitter
}

// NewCommandDispatcher creates new instance of CommandDispatcher
func NewCommandDispatcher(workspace *Workspace, transmitter services.Transmitter) CommandDispatcher {
	return CommandDispatcher{workspace: workspace, Transmitter: transmitter}
}

// DispatchCommand process incoming command by type and invoke specific method
func (dispatcher CommandDispatcher) DispatchCommand(message models.IncomingMessage) {
	user := dispatcher.workspace.users[message.Sender.NickName]
	switch message.Target {
	case CreateChannel:
		dispatcher.createChannel(user, message.Text)
	case JoinChannel:
		dispatcher.joinChannel(user, message.Text)
	case LeaveChannel:
		dispatcher.leaveChannel(user, message.Text)
	case Show:
		dispatcher.show(user, message.Text)
	}
}

func (dispatcher CommandDispatcher) createChannel(user entity.User, name string) {
	err := dispatcher.workspace.CreateChannel(user, name)
	if err != nil {
		dispatcher.SendMessageToUser(user, models.OutgoingMessage{Channel: name, Text: err.Error()})
		return
	}

	dispatcher.SendMessageToUser(user, dispatcher.workspaceChannelsMessage())
}

func (dispatcher CommandDispatcher) joinChannel(sender entity.User, name string) {
	user := dispatcher.workspace.users[sender.NickName]
	err := dispatcher.workspace.AddUserToChannel(sender, name)
	if err != nil {
		dispatcher.SendMessageToUser(user, models.OutgoingMessage{Channel: name, Text: err.Error()})
	}

	dispatcher.SendMessageToUser(user, dispatcher.channelSubscribersMessage(dispatcher.workspace.channels[name]))
}

func (dispatcher CommandDispatcher) leaveChannel(user entity.User, name string) {
	err := dispatcher.workspace.RemoveUserFromChannel(user, name)
	if err != nil {
		dispatcher.SendMessageToUser(user, models.OutgoingMessage{Channel: name, Text: err.Error()})
	}

	dispatcher.SendMessageToUser(user, dispatcher.channelSubscribersMessage(dispatcher.workspace.channels[name]))
}

func (dispatcher CommandDispatcher) show(user entity.User, param string) {
	switch param {
	case Users:
		dispatcher.SendMessageToUser(user, dispatcher.workspaceUsersMessage())
	case Channels:
		dispatcher.SendMessageToUser(user, dispatcher.workspaceChannelsMessage())
	case All:
		dispatcher.SendMessageToUser(user, dispatcher.workspaceChannelsMessage())
		dispatcher.SendMessageToUser(user, dispatcher.workspaceUsersMessage())
	default: // Get all users of specific channel if it exist
		if channel, exist := dispatcher.workspace.channels[param]; exist {
			dispatcher.SendMessageToUser(user, dispatcher.channelSubscribersMessage(channel))
		}
	}
}

func (dispatcher CommandDispatcher) workspaceChannelsMessage() models.OutgoingMessage {
	return models.OutgoingMessage{
		Channel: "Channels",
		Text:    strings.Join(dispatcher.workspace.WorkspaceChannels(), " | "),
		Time:    time.Now(),
	}
}

func (dispatcher CommandDispatcher) workspaceUsersMessage() models.OutgoingMessage {
	return models.OutgoingMessage{
		Sender: "Users",
		Text:   strings.Join(dispatcher.workspace.WorkspaceUsers(), " | "),
		Time:   time.Now(),
	}
}

func (dispatcher CommandDispatcher) channelSubscribersMessage(channel entity.Channel) models.OutgoingMessage {
	return models.OutgoingMessage{
		Channel: channel.Name,
		Sender:  "Users",
		Text:    strings.Join(dispatcher.workspace.ChannelSubscribers(channel), " | "),
		Time:    time.Now(),
	}
}
