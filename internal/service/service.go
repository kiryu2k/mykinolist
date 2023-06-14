package service

import "github.com/kiryu-dev/mykinolist/internal/model"

type AuthService interface {
	SignUp(user *model.User) error
	// SignIn(user *model.User) error
	// SignOut() error
	// Delete() error
}

type Service struct {
	AuthService
	// ListService
}

func New(repo AuthRepository) *Service {
	return &Service{
		AuthService: &authService{repo},
	}
}
