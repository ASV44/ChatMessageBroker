package main

import (
	"./components"
	"ChatMessageBroker/broker/entity"
	"ChatMessageBroker/broker/models"
	"encoding/json"
	"net"
	"time"
)

type Broker struct {
	server   *broker.Server
	users    []entity.User
	channels map[int]*broker.Channel
}

func main() {
	Broker := Broker{}
	Broker.Start()
}

func (Broker *Broker) Start() {
	Broker.server = &broker.Server{Host: broker.DEFAULT_HOST, Port: broker.DEFAULT_PORT, ConnectionType: broker.DEFAULT_TYPE}
	Broker.server.Start()
	defer Broker.server.Listener.Close()

	Broker.channels = make(map[int]*broker.Channel)
	Broker.channels[0] = &broker.Channel{Id: 0, Name: "random"}
	Broker.listen()

}

func (Broker *Broker) listen() {
	for {
		select {
		case connection := <-Broker.server.Connections:
			go Broker.register(connection)

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
}

//TODO: Create Channels(Rooms), Show all users and rooms of user at connecting to broker,
