package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kiryu-dev/mykinolist/internal/config"
	"github.com/kiryu-dev/mykinolist/internal/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenTTL  = 30 * time.Second // 30 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)

type authService struct {
	user  UserRepository
	token TokenRepository
	cfg   *config.Config
}

type UserRepository interface {
	CreateAccount(context.Context, *model.User) error
	FindByEmail(context.Context, string) (*model.User, error)
	UpdateLastLogin(context.Context, *model.User) error
	FindByID(context.Context, int64) (*model.User, error)
	DeleteAccount(context.Context, int64) error
}

type TokenRepository interface {
	Save(context.Context, *model.UserToken) error
	Remove(context.Context, string) error
}

func (s *authService) SignUp(userDTO *model.SignUpUserDTO) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := userDTO.Validate(); err != nil {
		return 0, err
	}
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.MinCost)
	if err != nil {
		return 0, err
	}
	user := &model.User{
		Username:       userDTO.Username,
		Email:          userDTO.Email,
		HashedPassword: string(HashedPassword),
		CreatedOn:      time.Now(),
		LastLogin:      time.Now(),
	}
	if err := s.user.CreateAccount(ctx, user); err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (s *authService) SignIn(userDTO *model.SignInUserDTO) (*model.Tokens, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := s.user.FindByEmail(ctx, userDTO.Email)
	if err != nil {
		return nil, fmt.Errorf("user with such email %s doesn't exist", userDTO.Email)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userDTO.Password))
	if err != nil {
		return nil, err
	}
	tokens, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, err
	}
	errChan := make(chan error, 2)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		user.LastLogin = time.Now()
		errChan <- s.user.UpdateLastLogin(ctx, user)
		wg.Done()
	}()
	go func() {
		errChan <- s.token.Save(ctx, &model.UserToken{
			UserID:       user.ID,
			RefreshToken: tokens.RefreshToken,
		})
		wg.Done()
	}()
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}
	return tokens, nil
}

func (s *authService) SignOut(refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.token.Remove(ctx, refreshToken)
}

func (s *authService) generateTokens(id int64) (*model.Tokens, error) {
	var (
		ATChan  = make(chan string, 1)
		RTChan  = make(chan string, 1)
		errChan = make(chan error, 2)
		wg      = &sync.WaitGroup{}
	)
	wg.Add(2)
	/* access token */
	go func() {
		ATPayload := &model.Payload{UserID: id, RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, ATPayload).
			SignedString([]byte(s.cfg.JWTAccessSecretKey))
		ATChan <- token
		errChan <- err
		wg.Done()
	}()
	/* refresh token */
	go func() {
		RTPayload := &model.Payload{UserID: id, RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, RTPayload).
			SignedString([]byte(s.cfg.JWTRefreshSecretKey))
		RTChan <- token
		errChan <- err
		wg.Done()
	}()
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			<-ATChan
			<-RTChan
			return nil, err
		}
	}
	return &model.Tokens{AccessToken: <-ATChan, RefreshToken: <-RTChan}, nil
}

func (s *authService) ParseAccessToken(tokenStr string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &model.Payload{}, func(token *jwt.Token) (any, error) {
		return []byte(s.cfg.JWTAccessSecretKey), nil
	})
	claims, ok := token.Claims.(*model.Payload)
	if !ok {
		return 0, err
	}
	if !token.Valid {
		return claims.UserID, &model.TokenError{Message: "token expiration date has passed"}
	}
	return claims.UserID, nil
}

func (s *authService) ParseRefreshToken(tokenStr string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &model.Payload{}, func(token *jwt.Token) (any, error) {
		return []byte(s.cfg.JWTRefreshSecretKey), nil
	})
	claims, ok := token.Claims.(*model.Payload)
	if !ok || !token.Valid {
		return 0, err
	}
	return claims.UserID, nil
}

func (s *authService) UpdateTokens(id int64) (*model.Tokens, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tokens, err := s.generateTokens(id)
	if err != nil {
		return nil, err
	}
	err = s.token.Save(ctx, &model.UserToken{
		UserID:       id,
		RefreshToken: tokens.RefreshToken,
	})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *authService) GetUser(id int64) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := s.user.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) Delete(id int64) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := s.user.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.user.DeleteAccount(ctx, id); err != nil {
		return nil, err
	}
	return user, nil
}
