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
	server            *broker.Server
	incoming          chan models.IncomingMessage
	users             map[string]entity.User
	channels          map[string]*broker.Channel
	messageDispatcher map[string]func(models.IncomingMessage)
	commandDispatcher map[string]func(models.User, string)
}

func main() {
	Broker := Broker{}
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
	text := "Welcome to Matrix workspace!\nEnter nickname:"
	message := models.Register{UserId: len(Broker.users), Text: text, Time: time.Now()}
	data, _ := json.Marshal(message)
	connection.Write(data)

	decoder := json.NewDecoder(connection)
	var user models.User
	decoder.Decode(&user)
	newUser := entity.User{Id: user.Id,
		NickName:   user.NickName,
		Connection: connection}

	Broker.users[newUser.NickName] = newUser
	subscribers := Broker.channels["random"].Subscribers
	Broker.channels["random"].Subscribers = append(subscribers, newUser)
	Broker.show(user, "all")

	go Broker.listen(connection)
	fmt.Printf("Connected user %s Id: %d\n", newUser.NickName, newUser.Id)
}

func (Broker *Broker) sendMessage(connection net.Conn, message models.OutcomingMessage) {
	data, _ := json.Marshal(message)
	connection.Write(data)
}

func (Broker *Broker) handleCommand(message models.IncomingMessage) {
	Broker.commandDispatcher[message.Target](message.Sender, message.Text)
}

func (Broker *Broker) handleDirectMessage(message models.IncomingMessage) {

}

func (Broker *Broker) handleChannelMessage(message models.IncomingMessage) {

}

func (Broker *Broker) createChannel(sender models.User, name string) {

}

func (Broker *Broker) joinChannel(sender models.User, name string) {

}

func (Broker *Broker) leaveChannel(sender models.User, name string) {

}

func (Broker *Broker) show(sender models.User, param string) {
	var channels []string
	var users []string

	for _, channel := range Broker.channels {
		channels = append(channels, channel.Name)
	}

	for _, user := range Broker.users {
		users = append(users, user.NickName)
	}

	channelsMessage := models.OutcomingMessage{Channel: "Channels", Text: strings.Join(channels, " "), Time: time.Now()}
	usersMessage := models.OutcomingMessage{Sender: "Users", Text: strings.Join(users, " "), Time: time.Now()}

	connection := Broker.users[sender.NickName].Connection

	switch param {
	case "users":
		Broker.sendMessage(connection, usersMessage)
	case "channels":
		Broker.sendMessage(connection, channelsMessage)
	case "all":
		Broker.sendMessage(connection, channelsMessage)
		if len(users) > 0 {
			Broker.sendMessage(connection, usersMessage)
		}
	}
}
