package main

import (
	"fmt"
	"github.com/ASV44/ChatMessageBroker/client/app"
	"os"
)

func main() {
	application := app.Init(app.DefaultHost, app.DefaultPort, app.DefaultType)

	err := application.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
