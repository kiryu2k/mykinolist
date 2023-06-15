package service

import (
	"time"

	"github.com/kiryu-dev/mykinolist/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo AuthRepository
}

type AuthRepository interface {
	CreateAccount(user *model.User) error
	LoginAccount(user *model.User) error
}

func (s *authService) SignUp(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	user.Password = string(encryptedPassword)
	user.CreatedOn = time.Now()
	user.LastLogin = user.CreatedOn
	if err := s.repo.CreateAccount(user); err != nil {
		return err
	}
	return nil
}
