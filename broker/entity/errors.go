package entity

import "fmt"

// AppInitFailed is returned when app initialization failed
type AppInitFailed struct {
	Message string
}

func (e AppInitFailed) Error() string {
	return fmt.Sprintf("Broker app init failed: %s", e.Message)
}

// ConfigInitFailed is returned when app failed to init config from config file at initialization phase
type ConfigInitFailed struct {
	Message string
}

func (e ConfigInitFailed) Error() string {
	return fmt.Sprintf("Init of config failed: %s", e.Message)
}

// WebsocketConfigDecodingFailed is returned when app failed to decode websocket config from map to struct
type WebsocketConfigDecodingFailed struct {
	Message string
}

func (e WebsocketConfigDecodingFailed) Error() string {
	return fmt.Sprintf("Decode of websocket config failed: %s", e.Message)
}

// UserWrongPassword is returned when new connection tries to sign in with existing username but wrong password is provided
type UserWrongPassword struct {
	Name string
}

func (e UserWrongPassword) Error() string {
	return fmt.Sprintf("Wrong password for existing user %s", e.Name)
}

// ChannelAlreadyExist is returned when user wants to create channel with name which already exist
type ChannelAlreadyExist struct {
	Name string
}

func (e ChannelAlreadyExist) Error() string {
	return fmt.Sprintf("%s channel already exist", e.Name)
}

// ChannelNotExist is returned when user wants to join channel which does not exist
type ChannelNotExist struct {
	Name string
}

func (e ChannelNotExist) Error() string {
	return fmt.Sprintf("Channel with name %s does not exits", e.Name)
}

// ChannelAlreadyJoined is returned when user wants to join channel which is already joined
type ChannelAlreadyJoined struct {
	Name string
}

func (e ChannelAlreadyJoined) Error() string {
	return fmt.Sprintf("%s channel already joined", e.Name)
}

// ChannelNotJoined is returned when user wants to leave channel which was not joined
type ChannelNotJoined struct {
	Name string
}

func (e ChannelNotJoined) Error() string {
	return fmt.Sprintf("not subscribed to %s!", e.Name)
}

// NotSupportedConnectionType is returned when socket connection type provided by launch argument is not supported
type NotSupportedConnectionType struct {
	ConnectionType string
}

func (e NotSupportedConnectionType) Error() string {
	return fmt.Sprintf("Socket connection type %s not supported!", e.ConnectionType)
}

// AuthServiceInitFailed is returned when broker fails to init auth service for managing authentication
type AuthServiceInitFailed struct {
	ErrorMessage string
}

func (e AuthServiceInitFailed) Error() string {
	return fmt.Sprintf("Failed to create Auth service: %s", e.ErrorMessage)
}

// TokenDecodingFailed is returned when JWT token has not expected signing method algorithm value encoded in token
type TokenDecodingFailed struct {
	Message string
}

func (e TokenDecodingFailed) Error() string {
	return fmt.Sprintf("Failed to decode auth token: %s", e.Message)
}

// InvalidToken is returned when received message from client contains invalid auth token
// which does not belong to user or is expired
type InvalidToken struct {
	Reason string
}

func (e InvalidToken) Error() string {
	return fmt.Sprintf("Received message with invalid auth token: %s", e.Reason)
}

// UserAuthFailed is returned when user authentication process failed at one step
type UserAuthFailed struct {
	Reason string
}

func (e UserAuthFailed) Error() string {
	return fmt.Sprintf("User authentication failed: %s", e.Reason)
}
