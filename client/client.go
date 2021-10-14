package client

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ASV44/chat-message-broker/client/models"

	"github.com/ASV44/chat-message-broker/client/components"

	"github.com/ASV44/chat-message-broker/client/models/receiver"
	"github.com/ASV44/chat-message-broker/client/models/sender"
)

// Client represents instance of client connection to broker
type Client struct {
	user        models.User
	inputReader components.InputReader
	commService components.CommunicationService
}

// NewClient creates new instance of Client app
func NewClient(user models.User, reader components.InputReader, commService components.CommunicationService) Client {
	return Client{
		user:        user,
		inputReader: reader,
		commService: commService,
	}
}

// Start init connection to broker and register new user on broker server
func (client Client) Start() {
	go client.listenConnection()
	client.listenUserInput()
}

func (client Client) listenConnection() {
	for {
		message, err := client.commService.GetMessage()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Connection closed by server ", err)
			} else {
				fmt.Println("Error at decoding message ", err)
			}

			return
		}

		client.showReceivedMessage(message)
	}
}

func (client Client) listenUserInput() {
	for {
		client.onUserAction(client.inputReader.GetUserInput())
	}
}

func (client Client) onUserAction(data string) {
	userInput := strings.Split(data, " ")
	operator := userInput[0][:1]
	target := userInput[0][1:]
	text := strings.Join(userInput[1:], " ")
	var messageType string

	switch operator {
	case "/":
		messageType = sender.CMD
	case "@":
		messageType = sender.DIRECT
	case "#":
		messageType = sender.CHANNEL
	}

	message := sender.Message{Type: messageType, Target: target, Sender: client.user, Text: text, Time: time.Now()}
	client.commService.SendMessage(message)
}

func (client Client) showReceivedMessage(message receiver.Message) {
	if message.Channel != "" {
		fmt.Printf("#%s ", message.Channel)
	}
	if message.Sender != "" {
		fmt.Printf("@%s ", message.Sender)
	}
	if message.Text != "" {
		fmt.Printf(": %s", message.Text)
	}
	if !message.Time.IsZero() {
		fmt.Printf("\t %v\n", message.Time)
	} else {
		fmt.Println()
	}
}
