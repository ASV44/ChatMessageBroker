package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ASV44/chat-message-broker/broker"
	"github.com/ASV44/chat-message-broker/broker/entity"
)

func main() {
	configFile := flag.String("config", "./broker/config.yaml", "Path to 'config.yaml' file")
	connection := flag.String(
		"connection",
		"all",
		"Type of broker socket connection. Available 'tcp-socket' and 'websocket'",
	)
	flag.Parse()
	brk, err := broker.Init(*configFile)
	if err != nil {
		log.Fatal(entity.AppInitFailed{Message: err.Error()})
	}
	err = brk.Start(*connection)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
