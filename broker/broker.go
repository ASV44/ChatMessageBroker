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
	server *broker.Server
	users  []entity.User
}

func main() {
	Broker := Broker{}
	Broker.Start()
}

func (Broker *Broker) Start() {
	Broker.server = &broker.Server{Host: broker.DEFAULT_HOST, Port: broker.DEFAULT_PORT, ConnectionType: broker.DEFAULT_TYPE}
	Broker.server.Start()
	defer Broker.server.Listener.Close()
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
	Broker.users = append(Broker.users, entity.User{Id: user.Id,
		NickName:   user.NickName,
		Connection: connection})
}

//TODO: Create Channels(Rooms), Show all users and rooms of user at connecting to broker,
