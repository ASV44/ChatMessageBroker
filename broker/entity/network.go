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

type Connection struct {
	NetworkConnection
	MessageIO
}

func NewRawTCPConnection(rawConnection net.Conn) Connection {
	return Connection{
		NetworkConnection: rawConnection,
		MessageIO:         NewJsonConnIO(rawConnection),
	}
}
