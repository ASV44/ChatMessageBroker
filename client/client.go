package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/ASV44/ChatMessageBroker/client/models"
	"github.com/ASV44/ChatMessageBroker/client/models/receiver"
	"github.com/ASV44/ChatMessageBroker/client/models/sender"
)

// Constant value of client config
const (
	DefaultHost = "localhost"
	DefaultPort = "8888"
	DefaultType = "tcp"
)

// Client represents instance of client connection to broker
type Client struct {
	connection  net.Conn
	user        models.User
	inputReader *bufio.Reader
	decoder     *json.Decoder
}

func main() {
	client := Client{}

	client.Start(DefaultType, DefaultHost, DefaultPort)
}

// Start init connection to broker and register new user on broker server
func (client Client) Start(connectionType string, host string, port string) {
	var err error
	fmt.Println("dialing")
	client.connection, err = net.Dial(connectionType, host+":"+port)
	fmt.Println("connecting")
	defer client.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client.decoder = json.NewDecoder(client.connection)
	client.inputReader = bufio.NewReader(os.Stdin)

	client.registerUser()

	go client.listenConnection()
	client.listenUserInput()
}

// Close end connection to broker
func (client Client) Close() {
	err := client.connection.Close()
	if err != nil {
		fmt.Println("Could not close client connection ", err)
	}
}

func (client Client) registerUser() {
	var registerMessage receiver.Register
	err := client.decoder.Decode(&registerMessage)
	if err != nil {
		fmt.Println("Could not decode register response ", err)
	}
	fmt.Println("Connected at: " + registerMessage.Time.Format("15:04:05 2006-01-02"))
	fmt.Print(registerMessage.Text)

	nickName := client.getUserInput()
	client.user = models.User{ID: registerMessage.UserId, NickName: nickName}
	userJSON, _ := json.Marshal(client.user)
	_, err = client.connection.Write(userJSON)
	if err != nil {
		fmt.Println("Could not write register data ", err)
	}
}

func (client Client) listenConnection() {
	var message receiver.Message
	for {
		if err := client.decoder.Decode(&message); err != io.EOF {
			client.showReceivedMessage(message)
		} else {
			return
		}
	}
}

func (client Client) listenUserInput() {
	for {
		client.onUserAction(client.getUserInput())
	}
}

func (client Client) getUserInput() string {
	data, _ := client.inputReader.ReadString('\n')
	return strings.TrimSuffix(string(data), "\n")
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
