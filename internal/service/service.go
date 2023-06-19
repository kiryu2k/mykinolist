package service

import (
	"github.com/kiryu-dev/mykinolist/internal/config"
	"github.com/kiryu-dev/mykinolist/internal/model"
)

type AuthService interface {
	SignUp(*model.SignUpUserDTO) (*model.User, error)
	SignIn(*model.SignInUserDTO) (*model.User, *model.Tokens, error)
	GetUser(int64) (*model.User, error)
	// SignOut() error
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
