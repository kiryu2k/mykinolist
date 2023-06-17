package service

import (
	"context"
	"fmt"
	"time"

	"github.com/kiryu-dev/mykinolist/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo AuthRepository
}

type AuthRepository interface {
	CreateAccount(context.Context, *model.User) error
	FindUserByEmail(context.Context, string) (*model.User, error)
	UpdateLastLogin(context.Context, *model.User) error
}

func (s *authService) SignUp(userDTO *model.SignUpUserDTO) (*model.User, error) {
	if err := userDTO.Validate(); err != nil {
		return nil, err
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Username:          userDTO.Username,
		Email:             userDTO.Email,
		EncryptedPassword: string(encryptedPassword),
		CreatedOn:         time.Now(),
		LastLogin:         time.Now(),
	}
	if err := s.repo.CreateAccount(context.Background(), user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) SignIn(userDTO *model.SignInUserDTO) (*model.User, error) {
	userFromDB, err := s.repo.FindUserByEmail(context.Background(), userDTO.Email)
	if err != nil {
		return nil, err
	}
	if userFromDB == nil {
		return nil, fmt.Errorf("user with such email %s doesn't exist", userDTO.Email)
	}
	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.EncryptedPassword), []byte(userDTO.Password))
	if err != nil {
		return nil, err
	}
	return userFromDB, nil
}
