package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kiryu-dev/mykinolist/internal/config"
	"github.com/kiryu-dev/mykinolist/internal/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenTTL  = 30 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)

type authService struct {
	user  UserRepository
	token TokenRepository
	cfg   *config.Config
}

type UserRepository interface {
	CreateAccount(context.Context, *model.User) error
	FindUserByEmail(context.Context, string) (*model.User, error)
	UpdateLastLogin(context.Context, *model.User) error
}

type TokenRepository interface {
	Save(context.Context, *model.UserToken) error
}

func (s *authService) SignUp(userDTO *model.SignUpUserDTO) (*model.User, *model.Tokens, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := userDTO.Validate(); err != nil {
		return nil, nil, err
	}
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.MinCost)
	if err != nil {
		return nil, nil, err
	}
	user := &model.User{
		Username:       userDTO.Username,
		Email:          userDTO.Email,
		HashedPassword: string(HashedPassword),
		CreatedOn:      time.Now(),
		LastLogin:      time.Now(),
	}
	if err := s.user.CreateAccount(ctx, user); err != nil {
		return nil, nil, err
	}
	tokens, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, nil, err
	}
	err = s.token.Save(ctx, &model.UserToken{
		UserID:       user.ID,
		RefreshToken: tokens.RefreshToken,
	})
	if err != nil {
		return nil, nil, err
	}
	return user, tokens, nil
}

func (s *authService) SignIn(userDTO *model.SignInUserDTO) (*model.User, *model.Tokens, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := s.user.FindUserByEmail(ctx, userDTO.Email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, fmt.Errorf("user with such email %s doesn't exist", userDTO.Email)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userDTO.Password))
	if err != nil {
		return nil, nil, err
	}
	tokens, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, nil, err
	}
	user.LastLogin = time.Now()
	if err := s.user.UpdateLastLogin(ctx, user); err != nil {
		return nil, nil, err
	}
	err = s.token.Save(ctx, &model.UserToken{
		UserID:       user.ID,
		RefreshToken: tokens.RefreshToken,
	})
	if err != nil {
		return nil, nil, err
	}
	return user, tokens, nil
}

func (s *authService) generateTokens(id int64) (*model.Tokens, error) {
	tokens := new(model.Tokens)
	/* access token payload */
	ATPayload := &model.Payload{UserID: id, RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}}
	var err error
	tokens.AccessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, ATPayload).
		SignedString([]byte(s.cfg.JWTAccessSecretKey))
	if err != nil {
		return nil, err
	}
	/* refresh token payload */
	RTPayload := &model.Payload{UserID: id, RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}}
	tokens.RefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, RTPayload).
		SignedString([]byte(s.cfg.JWTRefreshSecretKey))
	if err != nil {
		return nil, err
	}
	return tokens, nil
}
