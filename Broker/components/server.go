package broker

import (
	"fmt"
	"net"
	"os"
)

const (
	DEFAULT_HOST = "localhost"
	DEFAULT_PORT = "8888"
	DEFAULT_TYPE = "tcp"
)

type Server struct {
	Host           string
	Port           string
	ConnectionType string
	Listener       net.Listener
	Connections    chan net.Conn
}

func (server *Server) Start() {
	server.Connections = make(chan net.Conn)
	var err error
	server.Listener, err = net.Listen(server.ConnectionType, server.Host+":"+server.Port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("broker is running on port :", server.Port)
	}

	go server.acceptConnections()
}

func (server *Server) acceptConnections() {
	for {
		connection, err := server.Listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		server.Connections <- connection
	}
}
