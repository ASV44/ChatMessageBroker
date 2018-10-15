package entity

import "net"

type User struct {
	Id         int
	NickName   string
	Connection net.Conn
}
