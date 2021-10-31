package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHashing hashes password or checks hash password
type PasswordHashing interface {
	CompareHashAndPassword(passwordHash, plainPassword string) error
	HashPassword(password string) (string, error)
}

// LocalPasswordHashing hashes passwords with locally running BCrypt
type LocalPasswordHashing struct {
}

// NewLocalPasswordHashing creates new instance of local password hash
func NewLocalPasswordHashing() LocalPasswordHashing {
	return LocalPasswordHashing{}
}

// CompareHashAndPassword compares hash and password
func (e LocalPasswordHashing) CompareHashAndPassword(passwordHash, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plainPassword))

	if err != nil {
		return errors.New("hashed password is not the hash of the given password")
	}

	return nil
}

// HashPassword hashes a password and returns algorithm and hash string
func (e LocalPasswordHashing) HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordHash), nil
}
