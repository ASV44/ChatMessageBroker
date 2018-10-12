package brocker

type Broker struct {
	server *Server
}

func (broker *Broker) Init() *Broker {

	broker.server = &Server{host: DEFAULT_HOST, port: DEFAULT_PORT, connectionType: DEFAULT_TYPE}
	broker.server.Start()

	return broker
}
