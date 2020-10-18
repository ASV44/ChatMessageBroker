package entity

// User represents entity of broker user
type User struct {
	ID         int
	NickName   string
	Connection Connection
}
