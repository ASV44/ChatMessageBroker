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
	"time"
)

const (
	DEFAULT_HOST = "localhost"
	DEFAULT_PORT = "8888"
	DEFAULT_TYPE = "tcp"
)

type Client struct {
	connection net.Conn
	user       models.User
}

func main() {
	client := Client{}

	client.Start(DEFAULT_TYPE, DEFAULT_HOST, DEFAULT_PORT)
}

func (client *Client) Start(connectionType string, host string, port string) {
	var err error
	client.connection, err = net.Dial(connectionType, host+":"+port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go client.listen(bufio.NewReader(client.connection), client.onMessageReceive)
	client.listen(bufio.NewReader(os.Stdin), client.sendMessage)

}

func (client *Client) listen(reader *bufio.Reader, handler func(data []byte)) {
	for {
		data, _ := reader.ReadBytes('\n')
		fmt.Print("Receive from input: " + string(data))
		handler(data)
	}
}

func (client *Client) onMessageReceive(data []byte) {
	var message receiver.Message
	json.Unmarshal(data, &message)

}

func (client *Client) sendMessage(data []byte) {
	message := sender.Message{Sender: client.user, Text: string(data), Time: time.Now()}
	jsonData, _ := json.Marshal(message)
	client.connection.Write(jsonData)
}
