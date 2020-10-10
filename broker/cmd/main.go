package main

import (
	"github.com/ASV44/ChatMessageBroker/broker"
)

func main() {
	brk := broker.Init()
	brk.Start()
}
