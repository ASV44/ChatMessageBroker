package brocker

import "net"

type Broker struct {
	server *Server
}

func (broker *Broker) Start() {

	broker.server = &Server{host: DEFAULT_HOST, port: DEFAULT_PORT, connectionType: DEFAULT_TYPE}
	broker.server.Start()
	defer broker.server.listener.Close()
	broker.listen()

}

func (broker *Broker) listen() {
	for {
		select {
		case connection := <-broker.server.connections:
			go register(connection)

		}
	}
}

func register(connection net.Conn) {
	connection.Write([]byte("Welcome to Matrix chat!\n"))
	connection.Write([]byte("Select a nickname: "))
}
