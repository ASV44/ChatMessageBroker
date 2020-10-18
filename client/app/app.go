package app

import (
	"fmt"
	"github.com/ASV44/ChatMessageBroker/client"
	"github.com/ASV44/ChatMessageBroker/client/components"
	"github.com/ASV44/ChatMessageBroker/client/models"
	"github.com/ASV44/ChatMessageBroker/client/models/receiver"
	"github.com/ASV44/ChatMessageBroker/common"
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
	rawConnection, err := net.Dial(app.ConnectionType, app.Host+":"+app.Port)
	if err != nil {
		fmt.Println("Connection to server failed ", err)
		return err
	}

	connection := common.Connection{
		NetworkConnection: rawConnection,
		MessageIO:         common.NewJSONConnIO(rawConnection),
	}

	defer app.close(connection)

	user, err := app.registerUser(connection)
	if err != nil {
		fmt.Println("User registration failed ", err)
		return err
	}

	clientApp := client.NewClient(user, app.inputReader, components.NewCommunicationManager(connection))
	clientApp.Start()

	return nil
}

func (app App) registerUser(connection common.Connection) (models.User, error) {
	var registerMessage receiver.Register
	err := connection.GetMessage(&registerMessage)
	if err != nil {
		fmt.Println("Could not decode user register response ", err)
		return models.User{}, err
	}
	fmt.Println("Connected at: " + registerMessage.Time.Format("15:04:05 2006-01-02"))
	fmt.Print(registerMessage.Text)

	nickName := app.inputReader.GetUserInput()
	user := models.User{ID: registerMessage.UserID, NickName: nickName}
	err = connection.SendMessage(user)
	if err != nil {
		fmt.Println("Could not write user register data ", err)
		return user, err
	}

	return user, nil
}

// close end connection to broker
func (app App) close(connection common.Connection) {
	err := connection.Close()
	if err != nil {
		fmt.Println("Could not close client connection ", err)
	}
}
