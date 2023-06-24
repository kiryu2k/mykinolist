package model

type SearchResult struct {
	Docs []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"docs"`
}
