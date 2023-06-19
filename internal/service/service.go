package service

import (
	"github.com/kiryu-dev/mykinolist/internal/config"
	"github.com/kiryu-dev/mykinolist/internal/model"
)

type AuthService interface {
	SignUp(user *model.SignUpUserDTO) (*model.User, *model.Tokens, error)
	SignIn(user *model.SignInUserDTO) (*model.User, *model.Tokens, error)
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
