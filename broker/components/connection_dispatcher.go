package broker

import (
	"fmt"
	"github.com/ASV44/ChatMessageBroker/broker/entity"
	"github.com/ASV44/ChatMessageBroker/broker/models"
	"time"
)

type ConnectionManager interface {
	RegisterNewConnection(connection entity.Connection) error
}

type ConnectionDispatcher struct {
	workspace     *Workspace
	cmdDispatcher CommandDispatcher
}

func NewConnectionDispatcher(workspace *Workspace, cmdDispatcher CommandDispatcher) ConnectionDispatcher {
	return ConnectionDispatcher{
		workspace:     workspace,
		cmdDispatcher: cmdDispatcher,
	}
}

func (dispatcher ConnectionDispatcher) RegisterNewConnection(connection entity.Connection) error {
	text := fmt.Sprintf("Welcome to %s workspace!\nEnter nickname:", dispatcher.workspace.name)
	message := models.Register{UserId: len(dispatcher.workspace.users), Text: text, Time: time.Now()}
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

	fmt.Printf("Connected user: %s Id: %d addrr: %v\n", newUser.NickName, newUser.ID, connection.RemoteAddr())

	return nil
}
