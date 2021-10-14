package app

import (
	"fmt"
	"net"

	"github.com/ASV44/chat-message-broker/client"
	"github.com/ASV44/chat-message-broker/client/components"
	"github.com/ASV44/chat-message-broker/client/models"
	"github.com/ASV44/chat-message-broker/client/models/receiver"
	"github.com/ASV44/chat-message-broker/common"
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
		fmt.Println("Could not decode register message ", err)
		return models.User{}, err
	}

	fmt.Println("Connected at: " + registerMessage.Time.Format("15:04:05 2006-01-02"))
	fmt.Println(registerMessage.Text)

	err = app.registerUserInWorkspace(connection)
	if err != nil {
		fmt.Println("Workspace User registration failed", err)
		return models.User{}, err
	}

	var user models.User
	err = connection.GetMessage(&user)
	if err != nil {
		fmt.Println("Could not decode user register response ", err)
		return models.User{}, err
	}

	return user, nil
}

func (app App) registerUserInWorkspace(connection common.Connection) error {
	for {
		fmt.Print("Enter nickname: ")
		nickName := app.inputReader.GetUserInput()
		accountData := receiver.AccountData{NickName: nickName}
		err := connection.SendMessage(accountData)
		if err != nil {
			fmt.Println("Could not write user register data ", err)
			return err
		}

		var registrationMessage receiver.Message
		err = connection.GetMessage(&registrationMessage)
		if err != nil {
			fmt.Println("Could not decode registration message ", err)
			return err
		}

		switch registrationMessage.Sender {
		case "Error":
			fmt.Printf("%s: %s\n", registrationMessage.Sender, registrationMessage.Text)
		default:
			fmt.Println(registrationMessage.Text)
			return nil
		}
	}
}

// close end connection to broker
func (app App) close(connection common.Connection) {
	err := connection.Close()
	if err != nil {
		fmt.Println("Could not close client connection ", err)
	}
}
