package service

import (
	"context"
	"fmt"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

type listService struct {
	searcher MovieSearcher
	repo     MovieRepositroy
}

type MovieSearcher interface {
	Search(context.Context, string) (*model.SearchResult, error)
}

type MovieRepositroy interface {
	Add(context.Context, *model.Movie) error
}

func (s *listService) AddMovie(ctx context.Context, movie *model.Movie) error {
	if len(movie.Title) == 0 {
		return fmt.Errorf("empty movie name")
	}
	searchResult, err := s.searcher.Search(ctx, movie.Title)
	if err != nil {
		return err
	}
	movie.Title = searchResult.Docs[0].Name
	return s.repo.Add(ctx, movie)
}
