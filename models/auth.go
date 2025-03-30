package models

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
)

type AuthCredentials struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AuthRepository 数据访问接口
type AuthRepository interface {
	RegisterUser(ctx context.Context, registerData *AuthCredentials) (*User, error)
	GetUser(ctx context.Context, query interface{}, args ...interface{}) (*User, error)
}

// AuthService 业务逻辑接口
type AuthService interface {
	Login(ctx context.Context, loginData *AuthCredentials) (string, *User, error)
	Register(ctx context.Context, registerData *AuthCredentials) (string, *User, error)
}

// MatchesHash Check if a password matches a hash
func MatchesHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// IsValidEmail CHeck if an email is valid
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
