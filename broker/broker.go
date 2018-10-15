package main

import (
	"./components"
	"ChatMessageBroker/broker/models"
	"encoding/json"
	"net"
	"time"
)

type Broker struct {
	server *broker.Server
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
			go register(connection)

		}
	}
}

func register(connection net.Conn) {
	text := "Welcome to Matrix workspace!\nEnter nickname:"
	message := models.Message{Type: models.SYSTEM, Text: text, Time: time.Now()}
	data, _ := json.Marshal(message)
	connection.Write(append(data, '\n'))
}
