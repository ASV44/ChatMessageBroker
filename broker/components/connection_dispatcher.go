package broker

import (
	"fmt"
	"time"

	"github.com/ASV44/chat-message-broker/broker/entity"
	"github.com/ASV44/chat-message-broker/broker/models"
	"github.com/ASV44/chat-message-broker/broker/services"
	"github.com/ASV44/chat-message-broker/common"
)

// ConnectionManager represents abstraction of broker component which process new incoming connection
type ConnectionManager interface {
	RegisterNewConnection(connection common.Connection) (entity.User, error)
}

// ConnectionDispatcher represents broker component which process new incoming connection
type ConnectionDispatcher struct {
	workspace     *Workspace
	cmdDispatcher CommandDispatcher
	hashing       services.PasswordHashing
}

// NewConnectionDispatcher creates new instance of ConnectionDispatcher
func NewConnectionDispatcher(
	workspace *Workspace,
	cmdDispatcher CommandDispatcher,
	hashing services.PasswordHashing,
) ConnectionDispatcher {
	return ConnectionDispatcher{
		workspace:     workspace,
		cmdDispatcher: cmdDispatcher,
		hashing:       hashing,
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
	fmt.Println("###########", user, err)
	if err != nil {
		fmt.Println("Could not register new user in workspace", err)
		return user, err
	}

	err = dispatcher.sendSuccessfulAuthenticationMessage(connection)
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

		if user, ok := dispatcher.workspace.users[accountData.NickName]; ok {
			return dispatcher.signInExistingUserInWorkspace(user, accountData, connection)
		}

		passwordHash, err := dispatcher.hashing.HashPassword(accountData.Password)
		if err != nil {
			fmt.Println("User password hashing failed", err)
		}

		user, err := dispatcher.workspace.RegisterNewUser(
			entity.RegistrationData{
				NickName:     accountData.NickName,
				PasswordHash: passwordHash,
				Connection:   connection,
			},
		)
		switch err.(type) {
		case nil:
			return user, nil
		default:
			if err = dispatcher.sendRegistrationError(connection, err); err != nil {
				fmt.Println("Could not send registration error", err)
				return entity.User{}, err
			}
		}
	}
}

func (dispatcher ConnectionDispatcher) signInExistingUserInWorkspace(
	user entity.User,
	accountData models.AccountData,
	connection common.Connection,
) (entity.User, error) {
	if err := dispatcher.hashing.CompareHashAndPassword(user.PasswordHash, accountData.Password); err != nil {
		if err = dispatcher.sendRegistrationError(connection, err); err != nil {
			fmt.Println("Could not send registration error", err)
		}

		return entity.User{}, err
	}
	user.Connection = connection
	dispatcher.workspace.users[user.NickName] = user

	return user, nil
}

func (dispatcher ConnectionDispatcher) sendSuccessfulAuthenticationMessage(connection common.Connection) error {
	return connection.SendMessage(
		models.OutgoingMessage{
			Text: "Successfully authenticated in workspace",
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
