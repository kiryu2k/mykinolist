package service

import (
	"context"

	"github.com/kiryu-dev/mykinolist/internal/config"
	"github.com/kiryu-dev/mykinolist/internal/model"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type AuthService interface {
	SignUp(*model.SignUpUserDTO) (*model.ListInfo, error)
	SignIn(*model.SignInUserDTO) (*model.Tokens, error)
	SignOut(string) error
	GetUser(int64) (*model.User, error)
	ParseAccessToken(string) (int64, error)
	ParseRefreshToken(string) (int64, error)
	UpdateTokens(int64) (*model.Tokens, error)
	Delete(int64) (*model.User, error)
}

type ListService interface {
	AddMovie(context.Context, *model.ListUnit) error
	GetMovies(context.Context, int64) ([]*model.ListUnit, error)
	UpdateMovie(context.Context, *model.ListUnitPatch) error
}

type Service struct {
	AuthService
	ListService
}

func New(user UserRepository, token TokenRepository, list ListRepository,
	movie MovieRepositroy, searcher MovieSearcher, config *config.Config) *Service {
	return &Service{
		AuthService: &authService{user, token, list, config},
		ListService: &listService{searcher, movie, list},
	}
}
