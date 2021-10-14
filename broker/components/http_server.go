package broker

import (
	"fmt"
	"github.com/ASV44/chat-message-broker/broker/config"
	"github.com/gorilla/mux"
	"net"
	"net/http"
)

// HTTPServer represents instance of running http server
type HTTPServer struct {
	instance *http.Server
}

// InitHTTPServer creates and initialize instance of TCPServer
func InitHTTPServer(config config.Manager, router *mux.Router) HTTPServer {
	instance := &http.Server{
		Addr:         config.HTTPAddress(),
		Handler:      router,
		ReadTimeout:  config.ReadTimeout(),
		WriteTimeout: config.WriteTimeout(),
	}

	server := HTTPServer{
		instance: instance,
	}

	return server
}

// Start starts https server
func (server HTTPServer) Start() error {
	listener, err := net.Listen("tcp", server.instance.Addr)
	if err != nil {
		return err
	}

	go server.run(listener)

	fmt.Println("broker http server is running on :", server.instance.Addr)

	return nil
}

func (server HTTPServer) run(listener net.Listener) {
	err := server.instance.Serve(listener)
	if err != nil {
		fmt.Println("Error at starting to serve HTTP server ", err)
	}
}
