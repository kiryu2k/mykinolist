package service

import "github.com/kiryu-dev/mykinolist/internal/model"

type authService struct {
	repo AuthRepository
}

type AuthRepository interface {
	CreateAccount(user *model.User) error
	LoginAccount(user *model.User) error
}

func (s *authService) SignUp(user *model.User) error {
	return user.Validate()
}
