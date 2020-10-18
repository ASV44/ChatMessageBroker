package entity

import (
	"net"
)

type NetworkConnection interface {
	// Close closes the connection.
	Close() error

	// LocalAddr returns the local network address.
	LocalAddr() net.Addr

	// RemoteAddr returns the remote network address.
	RemoteAddr() net.Addr
}

type MessageIO interface {
	SendMessage(interface{}) error
	GetMessage(interface{}) error
}

type Connection struct {
	NetworkConnection
	MessageIO
}

func NewConnection(rawConnection net.Conn, io MessageIO) Connection {
	return Connection{
		NetworkConnection: rawConnection,
		MessageIO:         io,
	}
}
