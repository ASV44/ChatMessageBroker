package main

import (
	"./components"
	"net"
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
	connection.Write([]byte("Welcome to Matrix chat!\n"))
	connection.Write([]byte("Select a nickname:\n"))
}
