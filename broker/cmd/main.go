package main

import (
	"fmt"
	"github.com/ASV44/ChatMessageBroker/broker"
	"os"
)

func main() {
	brk := broker.Init()
	err := brk.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
