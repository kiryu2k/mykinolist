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
	SearchByID(context.Context, int64) (*model.Movie, error)
}

type MovieRepositroy interface {
	Add(context.Context, *model.ListUnit) error
	GetAll(context.Context, int64) ([]*model.ListUnit, error)
}

/* Add the first found movie by the specified title to the [kino]list */
func (s *listService) AddMovie(ctx context.Context, movie *model.ListUnit) error {
	if len(movie.Name) == 0 {
		return fmt.Errorf("empty movie name")
	}
	searchResult, err := s.searcher.Search(ctx, movie.Name)
	if err != nil {
		return err
	}
	movie.Movie = searchResult.Docs[0]
	return s.repo.Add(ctx, movie)
}

func (s *listService) GetMovies(ctx context.Context, userID int64) ([]*model.ListUnit, error) {
	movies, err := s.repo.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	for i, movieInfo := range movies {
		movie, err := s.searcher.SearchByID(ctx, movieInfo.ID)
		if err != nil {
			return nil, err
		}
		movies[i].Name = movie.Name
	}
	return movies, nil
}
