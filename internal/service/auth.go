package service

import (
	"fmt"
	"time"

	"github.com/kiryu-dev/mykinolist/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo AuthRepository
}

type AuthRepository interface {
	CreateAccount(*model.User) error
	FindUserByName(string) (*model.User, error)
	UpdateLastLogin(*model.User) error
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

func (s *authService) SignIn(user *model.User) error {
	userFromDB, err := s.repo.FindUserByName(user.Username)
	if err != nil {
		return err
	}
	if userFromDB == nil {
		return fmt.Errorf("user with name %s doesn't exist", user.Username)
	}
	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(user.Password))
	if err != nil {
		return err
	}
	user.ID = userFromDB.ID
	user.Password = userFromDB.Password
	user.CreatedOn = userFromDB.CreatedOn
	user.LastLogin = time.Now()
	return s.repo.UpdateLastLogin(user)
}
