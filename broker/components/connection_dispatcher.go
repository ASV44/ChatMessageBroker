package broker

import (
	"fmt"
	"github.com/ASV44/ChatMessageBroker/broker/entity"
	"github.com/ASV44/ChatMessageBroker/broker/models"
	"github.com/ASV44/ChatMessageBroker/common"
	"time"
)

// ConnectionManager represents abstraction of broker component which process new incoming connection
type ConnectionManager interface {
	RegisterNewConnection(connection common.Connection) (entity.User, error)
}

// ConnectionDispatcher represents broker component which process new incoming connection
type ConnectionDispatcher struct {
	workspace     *Workspace
	cmdDispatcher CommandDispatcher
}

// NewConnectionDispatcher creates new instance of ConnectionDispatcher
func NewConnectionDispatcher(workspace *Workspace, cmdDispatcher CommandDispatcher) ConnectionDispatcher {
	return ConnectionDispatcher{
		workspace:     workspace,
		cmdDispatcher: cmdDispatcher,
	}
}

// RegisterNewConnection register new incoming client connection by creating and adding new user to workspace
func (dispatcher ConnectionDispatcher) RegisterNewConnection(connection common.Connection) (entity.User, error) {
	text := fmt.Sprintf("Welcome to %s workspace!", dispatcher.workspace.name)
	message := models.Register{Text: text, Time: time.Now()}
	err := connection.SendMessage(message)
	if err != nil {
		fmt.Println("Could not send register data ", err)
		return entity.User{}, err
	}

	user, err := dispatcher.registerInWorkspace(connection)
	if err != nil {
		fmt.Println("Could not register new user in workspace", err)
		return user, err
	}

	err = dispatcher.sendSuccessfulRegistrationMessage(connection)
	if err != nil {
		fmt.Println("Could not send successful registration message ", err)
		dispatcher.workspace.RemoveUser(user)

		return user, err
	}

	userModel := models.User{ID: user.ID, NickName: user.NickName}
	err = connection.SendMessage(userModel)
	if err != nil {
		fmt.Println("Could not send user account data ", err)
		dispatcher.workspace.RemoveUser(user)

		return user, err
	}

	dispatcher.cmdDispatcher.show(user, All)

	fmt.Printf("Connected user: %s ID: %d addrr: %v\n", user.NickName, user.ID, connection.RemoteAddr())

	return user, nil
}

func (dispatcher ConnectionDispatcher) registerInWorkspace(connection common.Connection) (entity.User, error) {
	for {
		var accountData models.AccountData
		err := connection.GetMessage(&accountData)
		if err != nil {
			fmt.Println("Could not receive register data from user ", err)
			return entity.User{}, err
		}

		user, err := dispatcher.workspace.RegisterNewUser(
			entity.RegistrationData{NickName: accountData.NickName, Connection: connection},
		)
		switch err.(type) {
		case nil:
			return user, nil
		default:
			err = dispatcher.sendRegistrationError(connection, err)
			if err != nil {
				fmt.Println("Could not send registration error", err)
				return entity.User{}, err
			}
		}
	}
}

func (dispatcher ConnectionDispatcher) sendSuccessfulRegistrationMessage(connection common.Connection) error {
	return connection.SendMessage(
		models.OutgoingMessage{
			Text: "Successfully registered in workspace",
			Time: time.Now(),
		},
	)
}

func (dispatcher ConnectionDispatcher) sendRegistrationError(connection common.Connection, err error) error {
	return connection.SendMessage(
		models.OutgoingMessage{
			Sender: "Error",
			Text:   err.Error(),
			Time:   time.Now(),
		},
	)
}
