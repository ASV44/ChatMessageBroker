package main

import (
	"flag"
	"fmt"
	"github.com/ASV44/ChatMessageBroker/broker"
	"github.com/ASV44/ChatMessageBroker/broker/entity"

	"log"
	"os"
)

func main() {
	configFile := flag.String("config", "./broker/config.yaml", "Path to 'config.yaml' file")
	flag.Parse()
	brk, err := broker.Init(*configFile)
	if err != nil {
		log.Fatal(entity.AppInitFailed{Message: err.Error()})
	}
	err = brk.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
