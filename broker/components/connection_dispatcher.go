package broker

import (
	"fmt"
	"github.com/ASV44/ChatMessageBroker/broker/models"
	"github.com/ASV44/ChatMessageBroker/common"
	"time"
)

// ConnectionManager represents abstraction of broker component which process new incoming connection
type ConnectionManager interface {
	RegisterNewConnection(connection common.Connection) error
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
func (dispatcher ConnectionDispatcher) RegisterNewConnection(connection common.Connection) error {
	text := fmt.Sprintf("Welcome to %s workspace!\nEnter nickname:", dispatcher.workspace.name)
	message := models.Register{UserID: len(dispatcher.workspace.users), Text: text, Time: time.Now()}
	err := connection.SendMessage(message)
	if err != nil {
		fmt.Println("Could not send register data ", err)
		return err
	}

	var user models.User
	err = connection.GetMessage(&user)
	if err != nil {
		fmt.Println("Could not receive register data from user ", err)
		return err
	}
	newUser := user.ToUserEntity(connection)
	dispatcher.workspace.RegisterNewUser(newUser)
	dispatcher.cmdDispatcher.show(user, All)

	fmt.Printf("Connected user: %s ID: %d addrr: %v\n", newUser.NickName, newUser.ID, connection.RemoteAddr())

	return nil
}
