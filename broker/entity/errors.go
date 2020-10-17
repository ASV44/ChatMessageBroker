package entity

import "fmt"

type ChannelAlreadyExist struct {
	Name string
}

func (e ChannelAlreadyExist) Error() string {
	return fmt.Sprintf("%s channel already exist", e.Name)
}

type ChannelNotExist struct {
	Name string
}

func (e ChannelNotExist) Error() string {
	return fmt.Sprintf("Channel with name %s does not exits", e.Name)
}

type ChannelAlreadyJoined struct {
	Name string
}

func (e ChannelAlreadyJoined) Error() string {
	return fmt.Sprintf("%s channel already joined", e.Name)
}

type ChannelNotJoined struct {
	Name string
}

func (e ChannelNotJoined) Error() string {
	return fmt.Sprintf("not subscribed to %s!", e.Name)
}
