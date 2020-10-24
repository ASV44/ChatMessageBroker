package main

import (
	"flag"
	"fmt"
	"github.com/ASV44/ChatMessageBroker/client/app"
	"os"
)

func main() {
	host := flag.String("host", "0.0.0.0", "Broker host address")
	port := flag.String("port", "8888", "Broker host address port number")
	connectionType := flag.String("connection-type", "tcp", "Broker connection type")
	flag.Parse()

	application := app.Init(*host, *port, *connectionType)

	err := application.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
