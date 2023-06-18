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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := userDTO.Validate(); err != nil {
		return nil, err
	}
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Username:       userDTO.Username,
		Email:          userDTO.Email,
		HashedPassword: string(HashedPassword),
		CreatedOn:      time.Now(),
		LastLogin:      time.Now(),
	}
	if err := s.repo.CreateAccount(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) SignIn(userDTO *model.SignInUserDTO) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := s.repo.FindUserByEmail(ctx, userDTO.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user with such email %s doesn't exist", userDTO.Email)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userDTO.Password))
	if err != nil {
		return nil, err
	}
	user.LastLogin = time.Now()
	if err := s.repo.UpdateLastLogin(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
