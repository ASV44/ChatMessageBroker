package client

import (
	"encoding/json"
	"fmt"
	"github.com/ASV44/ChatMessageBroker/client/components"
	"io"
	"net"
	"strings"
	"time"

	"github.com/ASV44/ChatMessageBroker/client/models"
	"github.com/ASV44/ChatMessageBroker/client/models/receiver"
	"github.com/ASV44/ChatMessageBroker/client/models/sender"
)

// Client represents instance of client connection to broker
type Client struct {
	connection  net.Conn
	user        models.User
	decoder     *json.Decoder
	inputReader components.InputReader
}

func NewClient(connection net.Conn, user models.User, reader components.InputReader) Client {
	return Client{
		connection:  connection,
		user:        user,
		inputReader: reader,
		decoder:     json.NewDecoder(connection),
	}
}

// Start init connection to broker and register new user on broker server
func (client Client) Start() {
	go client.listenConnection()
	client.listenUserInput()
}

func (client Client) listenConnection() {
	var message receiver.Message
	for {
		if err := client.decoder.Decode(&message); err != io.EOF {
			client.showReceivedMessage(message)
		} else {
			fmt.Println("Error at decoding message from client connection ", err)
			return
		}
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

	client.sendMessage(messageType, target, text)
}

func (client Client) sendMessage(messageType string, target string, text string) {
	message := sender.Message{Type: messageType, Target: target, Sender: client.user, Text: text, Time: time.Now()}
	jsonData, _ := json.Marshal(message)
	_, err := client.connection.Write(jsonData)
	if err != nil {
		fmt.Println("Could not write client message ", err)
	}
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
