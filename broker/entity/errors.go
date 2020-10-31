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

// UserNameAlreadyExist is returned when new user tries to register with already registered nickname in workpsace
type UserNameAlreadyExist struct {
	Name string
}

func (e UserNameAlreadyExist) Error() string {
	return fmt.Sprintf("%s user name already exist", e.Name)
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
