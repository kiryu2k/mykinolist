package service

import (
	"github.com/kiryu-dev/mykinolist/internal/config"
	"github.com/kiryu-dev/mykinolist/internal/model"
)

type AuthService interface {
	SignUp(*model.SignUpUserDTO) (int64, error)
	SignIn(*model.SignInUserDTO) (*model.Tokens, error)
	SignOut(string) error
	GetUser(int64) (*model.User, error)
	ParseAccessToken(string) (int64, error)
	ParseRefreshToken(string) (int64, error)
	UpdateTokens(int64) (*model.Tokens, error)
	// Delete() error
}

type Service struct {
	AuthService
	// ListService
}

func New(user UserRepository, token TokenRepository, config *config.Config) *Service {
	return &Service{
		AuthService: &authService{user, token, config},
	}
}
