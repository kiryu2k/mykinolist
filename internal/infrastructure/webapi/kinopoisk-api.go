package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

const (
	urlApiSearchRequest     = "https://api.kinopoisk.dev/v1.2/movie/search?page=1&limit=1&query=%s"
	urlApiSearchByIDRequest = "https://api.kinopoisk.dev/v1.3/movie/%d"
)

type KinopoiskWebAPI struct {
	apiKey string
}

func New(apiKey string) *KinopoiskWebAPI {
	return &KinopoiskWebAPI{apiKey: apiKey}
}

func (api *KinopoiskWebAPI) Search(ctx context.Context, title string) (*model.SearchResult, error) {
	url := fmt.Sprintf(urlApiSearchRequest, title)
	searchResult := new(model.SearchResult)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-KEY", api.apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(searchResult); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return searchResult, nil
}

func (api *KinopoiskWebAPI) SearchByID(ctx context.Context, id int64) (*model.Movie, error) {
	url := fmt.Sprintf(urlApiSearchByIDRequest, id)
	movie := new(model.Movie)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-KEY", api.apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(movie); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return movie, nil
}
