package entity

import "net"

// User represents entity of broker user
type User struct {
	ID         int
	NickName   string
	Connection net.Conn
}
