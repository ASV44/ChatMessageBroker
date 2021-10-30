package services

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	bitSize = 4096
)

type AuthProvider interface {
	GenerateNewJWT(claims *jwt.RegisteredClaims) (string, error)
}

type AuthService struct {
	privateKey *rsa.PrivateKey
}

func NewAuthService() (AuthService, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return AuthService{}, err
	}

	return AuthService{
		privateKey: privateKey,
	}, nil
}

func (a AuthService) GenerateNewJWT(claims *jwt.RegisteredClaims) (string, error) {
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = uuid.NewString()

	// Generate encoded token
	return token.SignedString(a.privateKey)
}
