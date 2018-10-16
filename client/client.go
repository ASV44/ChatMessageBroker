package main

import (
	"ChatMessageBroker/client/models"
	"ChatMessageBroker/client/models/receiver"
	"ChatMessageBroker/client/models/sender"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	DEFAULT_HOST = "localhost"
	DEFAULT_PORT = "8888"
	DEFAULT_TYPE = "tcp"
)

type Client struct {
	connection  net.Conn
	user        models.User
	inputReader *bufio.Reader
	decoder     *json.Decoder
}

func main() {
	client := Client{}

	client.Start(DEFAULT_TYPE, DEFAULT_HOST, DEFAULT_PORT)
}

func (client *Client) Start(connectionType string, host string, port string) {
	var err error
	client.connection, err = net.Dial(connectionType, host+":"+port)
	defer client.connection.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go client.sigIntHandler()

	client.decoder = json.NewDecoder(client.connection)
	client.inputReader = bufio.NewReader(os.Stdin)

	client.registerUser()

	go client.listenConnection()
	client.listenUserInput()
}

func (client *Client) sigIntHandler() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
	client.connection.Close()
	os.Exit(0)
}

func (client *Client) registerUser() {
	var registerMessage receiver.Register
	client.decoder.Decode(&registerMessage)
	fmt.Println("Connected at: " + registerMessage.Time.Format("15:04:05 2006-01-02"))
	fmt.Print(registerMessage.Text)

	nickName := client.getUserInput()
	client.user = models.User{Id: registerMessage.UserId, NickName: nickName}
	userJson, _ := json.Marshal(client.user)
	client.connection.Write(userJson)
}

func (client *Client) listenConnection() {
	var message receiver.Message
	for {
		client.decoder.Decode(&message)
	}
}

func (client *Client) listenUserInput() {
	for {
		client.onUserAction(client.getUserInput())
	}
}

func (client *Client) getUserInput() string {
	data, _ := client.inputReader.ReadString('\n')
	return strings.TrimSuffix(string(data), "\n")
}

func (client *Client) onMessageReceive(data []byte) {
	var message receiver.Message
	json.Unmarshal(data, &message)

}

func (client *Client) onUserAction(data string) {
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

func (client *Client) sendMessage(messageType string, target string, text string) {
	message := sender.Message{Type: messageType, Target: target, Sender: client.user, Text: text, Time: time.Now()}
	jsonData, _ := json.Marshal(message)
	client.connection.Write(jsonData)
}

//TODO: Implement commands /open {room} /start {name} /create {room}
