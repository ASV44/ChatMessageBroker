package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	DEFAULT_HOST = "localhost"
	DEFAULT_PORT = "8888"
	DEFAULT_TYPE = "tcp"
)

type Client struct {
	connection net.Conn
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

	for {
		//reader := bufio.NewReader(os.Stdin)
		//text, _ := reader.ReadString('\n')
		//fmt.Fprintf(conn, text + "\n")
		//// listen for reply
		message, _ := bufio.NewReader(client.connection).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}

}
