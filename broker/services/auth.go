package services

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/ASV44/chat-message-broker/broker/entity"
)

const (
	bitSize  = 4096
	KeyIDKey = "kid"
)

// AuthProvider defines operations of manging clients authentication
type AuthProvider interface {
	GenerateNewJWT(*jwt.RegisteredClaims) (string, error)
	GenerateNewUserJWT(userID string, clientAddr string) (string, error)
	DecodeToken(string) (*jwt.Token, jwt.RegisteredClaims, error)
}

// AuthService implements AuthProvider and manages clients authentication
type AuthService struct {
	tokenIssuer string
	privateKey  *rsa.PrivateKey
	keyID       string
}

// NewAuthService creates new instance of AuthService
func NewAuthService(tokenIssuer string) (AuthService, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return AuthService{}, err
	}

	return AuthService{
		tokenIssuer: tokenIssuer,
		privateKey:  privateKey,
		keyID:       uuid.NewString(),
	}, nil
}

// GenerateNewJWT creates new JWT token signed by RSA256 private key
func (a AuthService) GenerateNewJWT(claims *jwt.RegisteredClaims) (string, error) {
	claims.Issuer = a.tokenIssuer
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header[KeyIDKey] = a.keyID

	// Generate encoded token
	return token.SignedString(a.privateKey)
}

// GenerateNewUserJWT creates new JWT token with dedicated user specific configuration
func (a AuthService) GenerateNewUserJWT(userID string, clientAddr string) (string, error) {
	return a.GenerateNewJWT(&jwt.RegisteredClaims{
		Subject:   userID,
		Audience:  jwt.ClaimStrings{fmt.Sprintf("client %s", clientAddr)},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
}

// DecodeToken parse provided JWT token string and returns token instance and token claims
func (a AuthService) DecodeToken(tokenString string) (*jwt.Token, jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, entity.TokenDecodingFailed{
				Message: fmt.Sprintf(
					"Unexpected signing method. Used signing %s is not RS256",
					token.Header["alg"].(string),
				),
			}
		}

		if token.Header[KeyIDKey] != a.keyID {
			return nil, entity.TokenDecodingFailed{Message: "Token key ID missing or not equal to expected"}
		}

		return &a.privateKey.PublicKey, nil
	})

	return token, claims, err
}
