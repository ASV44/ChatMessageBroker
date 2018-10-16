package main

import (
	"./components"
	"ChatMessageBroker/broker/entity"
	"ChatMessageBroker/broker/models"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

type Broker struct {
	server            *broker.Server
	users             []entity.User
	incoming          chan models.IncomingMessage
	channels          map[int]*broker.Channel
	messageDispatcher map[string]func(message models.IncomingMessage)
	commandDispatcher map[string]func(param string)
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
	Broker.incoming = make(chan models.IncomingMessage)
	Broker.channels = make(map[int]*broker.Channel)
	Broker.channels[0] = &broker.Channel{Id: 0, Name: "random"}
	Broker.run()

}

func (Broker *Broker) initMessageDispatcher() {
	Broker.messageDispatcher = make(map[string]func(message models.IncomingMessage))

	Broker.messageDispatcher[models.CMD] = Broker.handleCommand
	Broker.messageDispatcher[models.DIRECT] = Broker.handleDirectMessage
	Broker.messageDispatcher[models.CHANNEL] = Broker.handleChannelMessage
}

func (Broker *Broker) initCommandDispatcher() {
	Broker.commandDispatcher = make(map[string]func(param string))

	Broker.commandDispatcher["create"] = Broker.handleCommand
	Broker.messageDispatcher[models.DIRECT] = Broker.handleDirectMessage
	Broker.messageDispatcher[models.CHANNEL] = Broker.handleChannelMessage
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
			go Broker.handleIncomingMessages(message)
		}
	}
}

func (Broker *Broker) register(connection net.Conn) {
	text := "Welcome to Matrix workspace!\nEnter nickname:"
	message := models.Register{UserId: len(Broker.users), Text: text, Time: time.Now()}
	data, _ := json.Marshal(message)
	connection.Write(append(data, '\n'))

	decoder := json.NewDecoder(connection)
	var user models.User
	decoder.Decode(&user)
	newUser := entity.User{Id: user.Id,
		NickName:   user.NickName,
		Connection: connection}

	Broker.users = append(Broker.users, newUser)
	Broker.channels[0].Subscribers = Broker.users

	go Broker.listen(connection)
	fmt.Printf("Connected user %s Id: %d\n", newUser.NickName, newUser.Id)
}

func (Broker *Broker) handleIncomingMessages(message models.IncomingMessage) {
	fmt.Println(message.Type, message.Target, message.Sender, message.Text, message.Time)
}

func (Broker *Broker) handleCommand(message models.IncomingMessage) {

}

func (Broker *Broker) handleDirectMessage(message models.IncomingMessage) {

}

func (Broker *Broker) handleChannelMessage(message models.IncomingMessage) {

}

func (Broker *Broker) createChannel(name string) {

}

func (Broker *Broker) joinChannel(name string) {

}

func (Broker *Broker) leaveChannel(name string) {

}

func (Broker *Broker) show(param string) {

}

//TODO: Create Channels(Rooms), Show all users and rooms of user at connecting to broker,
