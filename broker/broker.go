package main

import (
	"./components"
	"ChatMessageBroker/broker/entity"
	"ChatMessageBroker/broker/models"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type Broker struct {
	workspace         string
	server            *broker.Server
	incoming          chan models.IncomingMessage
	users             map[string]entity.User
	channels          map[string]*broker.Channel
	messageDispatcher map[string]func(models.IncomingMessage)
	commandDispatcher map[string]func(models.User, string)
}

func main() {
	Broker := Broker{workspace: "Matrix"}
	Broker.Start()
}

func (Broker *Broker) Start() {
	Broker.server = &broker.Server{Host: broker.DEFAULT_HOST, Port: broker.DEFAULT_PORT, ConnectionType: broker.DEFAULT_TYPE}
	Broker.server.Start()
	defer Broker.server.Listener.Close()

	Broker.initMessageDispatcher()
	Broker.initCommandDispatcher()
	Broker.incoming = make(chan models.IncomingMessage)
	Broker.users = make(map[string]entity.User)
	Broker.channels = make(map[string]*broker.Channel)
	Broker.channels["random"] = &broker.Channel{Id: 0, Name: "random"}
	Broker.run()

}

func (Broker *Broker) initMessageDispatcher() {
	Broker.messageDispatcher = make(map[string]func(models.IncomingMessage))

	Broker.messageDispatcher[models.CMD] = Broker.handleCommand
	Broker.messageDispatcher[models.DIRECT] = Broker.handleDirectMessage
	Broker.messageDispatcher[models.CHANNEL] = Broker.handleChannelMessage
}

func (Broker *Broker) initCommandDispatcher() {
	Broker.commandDispatcher = make(map[string]func(models.User, string))

	Broker.commandDispatcher["create"] = Broker.createChannel
	Broker.commandDispatcher["join"] = Broker.joinChannel
	Broker.commandDispatcher["leave"] = Broker.leaveChannel
	Broker.commandDispatcher["show"] = Broker.show
}

func (Broker *Broker) listen(connection net.Conn) {
	decoder := json.NewDecoder(connection)
	var message models.IncomingMessage
	for {
		if err := decoder.Decode(&message); err != io.EOF {
			Broker.incoming <- message
		} else {
			return
		}
	}
}

func (Broker *Broker) run() {
	for {
		select {
		case connection := <-Broker.server.Connections:
			go Broker.register(connection)
		case message := <-Broker.incoming:
			go Broker.messageDispatcher[message.Type](message)
		}
	}
}

func (Broker *Broker) register(connection net.Conn) {
	text := fmt.Sprintf("Welcome to %s workspace!\nEnter nickname:", Broker.workspace)
	message := models.Register{UserId: len(Broker.users), Text: text, Time: time.Now()}
	data, _ := json.Marshal(message)
	connection.Write(data)

	decoder := json.NewDecoder(connection)
	var user models.User
	decoder.Decode(&user)
	newUser := user.ToUserEntity(connection)

	Broker.users[newUser.NickName] = newUser
	subscribers := Broker.channels["random"].Subscribers
	Broker.channels["random"].Subscribers = append(subscribers, newUser)
	Broker.show(user, "all")

	go Broker.listen(connection)
	fmt.Printf("Connected user: %s Id: %d addrr: %v\n", newUser.NickName, newUser.Id, connection.RemoteAddr())
}

func (Broker *Broker) sendMessage(connection net.Conn, message models.OutcomingMessage) {
	data, _ := json.Marshal(message)
	connection.Write(data)
}

func (Broker *Broker) handleCommand(message models.IncomingMessage) {
	Broker.commandDispatcher[message.Target](message.Sender, message.Text)
}

func (Broker *Broker) handleDirectMessage(message models.IncomingMessage) {
	user, isPresent := Broker.users[message.Target]
	if isPresent {
		Broker.sendMessage(user.Connection, message.ToOutcomingMessage())
	}
}

func (Broker *Broker) handleChannelMessage(message models.IncomingMessage) {
	channel, isPresent := Broker.channels[message.Target]
	if isPresent {
		sender := Broker.users[message.Sender.NickName]
		if channel.Contains(sender) {
			draft := message.ToOutcomingMessage()
			for _, user := range channel.Subscribers {
				if user.Id != message.Sender.Id {
					go Broker.sendMessage(user.Connection, draft)
				}
			}
		} else {
			Broker.sendMessage(sender.Connection, models.OutcomingMessage{Channel: channel.Name, Text: "not joined yet!"})
		}
	}
}

func (Broker *Broker) createChannel(sender models.User, name string) {
	user := Broker.users[sender.NickName]
	_, exist := Broker.channels[name]
	if !exist {
		channel := &broker.Channel{Id: len(Broker.channels), Name: name}
		Broker.channels[name] = channel
		channel.Subscribers = append(channel.Subscribers, user)
		Broker.sendMessage(user.Connection, Broker.getWorkspaceChannelsMessage())
	} else {
		Broker.sendMessage(user.Connection, models.OutcomingMessage{Channel: name, Text: "already exist!"})
	}
}

func (Broker *Broker) joinChannel(sender models.User, name string) {
	user := Broker.users[sender.NickName]
	channel, exist := Broker.channels[name]
	if exist {
		if !channel.Contains(user) {
			channel.Subscribers = append(channel.Subscribers, user)
			Broker.sendMessage(user.Connection, Broker.getChannelSubscribersMessage(channel))
		} else {
			Broker.sendMessage(user.Connection, models.OutcomingMessage{Channel: name, Text: "already joined!"})
		}
	} else {
		Broker.sendMessage(user.Connection, models.OutcomingMessage{Channel: name, Text: "does not exist!"})
	}
}

func (Broker *Broker) leaveChannel(sender models.User, name string) {
	user := Broker.users[sender.NickName]
	channel, exist := Broker.channels[name]
	if exist {
		if isPresent, index := channel.ContainsSubscriber(user); isPresent {
			channel.Subscribers = append(channel.Subscribers[:index], channel.Subscribers[index+1:]...)
			Broker.sendMessage(user.Connection, Broker.getChannelSubscribersMessage(channel))
		} else {
			Broker.sendMessage(user.Connection, models.OutcomingMessage{Channel: name, Text: "not subscribed to!"})
		}
	} else {
		Broker.sendMessage(user.Connection, models.OutcomingMessage{Channel: name, Text: "does not exist!"})
	}
}

func (Broker *Broker) getWorkspaceChannelsMessage() models.OutcomingMessage {
	var channels []string
	for _, channel := range Broker.channels {
		channels = append(channels, channel.Name)
	}

	return models.OutcomingMessage{Channel: "Channels",
		Text: strings.Join(channels, " | "),
		Time: time.Now()}
}

func (Broker *Broker) getWorkspaceUsersMessage() models.OutcomingMessage {
	var users []string
	for _, user := range Broker.users {
		users = append(users, user.NickName)
	}
	return models.OutcomingMessage{Sender: "Users",
		Text: strings.Join(users, " | "),
		Time: time.Now()}

}

func (Broker *Broker) getChannelSubscribersMessage(channel *broker.Channel) models.OutcomingMessage {
	var users []string
	for _, user := range channel.Subscribers {
		users = append(users, user.NickName)
	}
	return models.OutcomingMessage{Channel: channel.Name,
		Sender: "Users",
		Text:   strings.Join(users, " | "),
		Time:   time.Now()}
}

func (Broker *Broker) show(sender models.User, param string) {
	connection := Broker.users[sender.NickName].Connection
	switch param {
	case "users":
		Broker.sendMessage(connection, Broker.getWorkspaceUsersMessage())
	case "channels":
		Broker.sendMessage(connection, Broker.getWorkspaceChannelsMessage())
	case "all":
		Broker.sendMessage(connection, Broker.getWorkspaceChannelsMessage())
		Broker.sendMessage(connection, Broker.getWorkspaceUsersMessage())
	default:
		channel, exist := Broker.channels[param]
		if exist {
			Broker.sendMessage(connection, Broker.getChannelSubscribersMessage(channel))
		}
	}
}
