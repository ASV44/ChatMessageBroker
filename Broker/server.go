package brocker

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
	host           string
	port           string
	connectionType string
	listener       net.Listener
	connections    chan net.Conn
}

func (server *Server) Start() {
	server.connections = make(chan net.Conn)
	var err error
	server.listener, err = net.Listen(server.connectionType, server.host+":"+server.port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Broker is running on port :", server.port)
	}
	defer server.listener.Close()

	go server.acceptConnections()
}

func (server *Server) acceptConnections() {
	for {
		connection, err := server.listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		server.connections <- connection
	}
}
