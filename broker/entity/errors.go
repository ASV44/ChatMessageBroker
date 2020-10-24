package entity

import "fmt"

// ConfigInitFailed is returned when app failed to init config from config file at initialization phase
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
