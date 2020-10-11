package app

import (
	"encoding/json"
	"fmt"
	"github.com/ASV44/ChatMessageBroker/client"
	"github.com/ASV44/ChatMessageBroker/client/components"
	"github.com/ASV44/ChatMessageBroker/client/models"
	"github.com/ASV44/ChatMessageBroker/client/models/receiver"
	"net"
)

// Constant value of client config
const (
	DefaultHost = "localhost"
	DefaultPort = "8888"
	DefaultType = "tcp"
)

// App represents instance of the client app
type App struct {
	Host           string
	Port           string
	ConnectionType string
	inputReader    components.InputReader
}

// Init creates and initialize new instance of App
func Init(host string, port string, connectionType string) App {
	return App{
		Host:           host,
		Port:           port,
		ConnectionType: connectionType,
		inputReader:    components.NewInputReader(),
	}
}

// Start init connection to broker and register new user on broker server
func (app App) Start() error {
	connection, err := net.Dial(app.ConnectionType, app.Host+":"+app.Port)
	if err != nil {
		fmt.Println("Connection to server failed ", err)
		return err
	}
	defer app.close(connection)

	decoder := json.NewDecoder(connection)

	user, err := app.registerUser(connection, decoder)
	if err != nil {
		fmt.Println("User registration failed ", err)
		return err
	}

	clientApp := client.NewClient(connection, user, app.inputReader)
	clientApp.Start()

	return nil
}

func (app App) registerUser(connection net.Conn, decoder *json.Decoder) (models.User, error) {
	var registerMessage receiver.Register
	err := decoder.Decode(&registerMessage)
	if err != nil {
		fmt.Println("Could not decode user register response ", err)
		return models.User{}, err
	}
	fmt.Println("Connected at: " + registerMessage.Time.Format("15:04:05 2006-01-02"))
	fmt.Print(registerMessage.Text)

	nickName := app.inputReader.GetUserInput()
	user := models.User{ID: registerMessage.UserId, NickName: nickName}
	userJSON, _ := json.Marshal(user)
	_, err = connection.Write(userJSON)
	if err != nil {
		fmt.Println("Could not write user register data ", err)
		return user, err
	}

	return user, nil
}

// close end connection to broker
func (app App) close(connection net.Conn) {
	err := connection.Close()
	if err != nil {
		fmt.Println("Could not close client connection ", err)
	}
}
