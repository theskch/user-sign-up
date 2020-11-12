package authentication

import (
	"fmt"
	"time"

	"github.com/dchest/passwordreset"
	"golang.org/x/crypto/bcrypt"
)

// PasswordResetHandlerInterface used for all password related operations
type PasswordResetHandlerInterface interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password string, hash string) bool

	GenerateResetPasswordToken(login string, hashedPassword string) (token string, passwordResetLink string)
	VerifyResetPasswordToken(token string, hashedPasswordFunc func(string) ([]byte, error)) (login string, err error)
}

// PasswordResetHandler implementation
type PasswordResetHandler struct {
	Settings PasswordResetSettings
}

// HashPassword is used to hide user password
func (handler *PasswordResetHandler) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost+10)
	return string(bytes), err
}

// CheckPasswordHash checks if password hash matches the provided password
func (handler *PasswordResetHandler) CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateResetPasswordToken generates password reset token and returnes password reset url
func (handler *PasswordResetHandler) GenerateResetPasswordToken(login string, hashedPassword string) (token string, passwordResetLink string) {
	token = passwordreset.NewToken(login, time.Duration(handler.Settings.TokenExpirationTime)*time.Minute,
		[]byte(hashedPassword), []byte(handler.Settings.TokenSecretKey))
	passwordResetLink = fmt.Sprintf("%s/resetpass?token=%s", handler.Settings.WebsiteBaseURL, token)
	return
}

// VerifyResetPasswordToken verifies password reset token
func (handler *PasswordResetHandler) VerifyResetPasswordToken(token string, hashedPasswordFunc func(string) ([]byte, error)) (login string, err error) {
	return passwordreset.VerifyToken(
		token,
		hashedPasswordFunc,
		[]byte(handler.Settings.TokenSecretKey))
}
