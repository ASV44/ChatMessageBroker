package common

import (
	"net"
)

// NetworkConnection represents abstraction of network connection
type NetworkConnection interface {
	// Close closes the connection.
	Close() error

	// LocalAddr returns the local network address.
	LocalAddr() net.Addr

	// RemoteAddr returns the remote network address.
	RemoteAddr() net.Addr
}

// MessageIO represents abstraction of message input output of network connection
type MessageIO interface {
	SendMessage(interface{}) error
	GetMessage(interface{}) error
}

// Connection represents entity for abstraction of broker connection
type Connection struct {
	NetworkConnection
	MessageIO
}

// NewConnection creates new instance of Connection
func NewConnection(rawConnection NetworkConnection, io MessageIO) Connection {
	return Connection{
		NetworkConnection: rawConnection,
		MessageIO:         io,
	}
}
