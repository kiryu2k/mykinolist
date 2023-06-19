package service

import (
	"github.com/kiryu-dev/mykinolist/internal/config"
	"github.com/kiryu-dev/mykinolist/internal/model"
)

type AuthService interface {
	SignUp(user *model.SignUpUserDTO) (*model.User, error)
	SignIn(user *model.SignInUserDTO) (*model.Tokens, error)
	// SignOut() error
	// Delete() error
}

type Service struct {
	AuthService
	// ListService
}

func New(repo UserRepository, config *config.Config) *Service {
	return &Service{
		AuthService: &authService{repo, config},
	}
}
