package authentication

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTHandlerInterface used for handling jwt tokens
type JWTHandlerInterface interface {
	GenerateToken(userID string) (string, error)
	ParseToken(jwtString string, claims jwt.Claims) (*jwt.Token, error)
}

// JWTHandler implemntation
type JWTHandler struct {
	Settings JWTSettings
}

// GenerateToken implemntation
func (handler *JWTHandler) GenerateToken(userID string) (token string, err error) {
	tokenBuilder := jwt.New(jwt.SigningMethodHS512)

	tokenBuilder.Claims =
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(handler.Settings.JWTExpirationTime)).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   userID,
		}
	return tokenBuilder.SignedString([]byte(handler.Settings.JWTSecretKey))
}

// ParseToken implementation
func (handler *JWTHandler) ParseToken(jwtString string, claims jwt.Claims) (token *jwt.Token, err error) {
	return jwt.ParseWithClaims(jwtString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(handler.Settings.JWTSecretKey), nil
	})
}
