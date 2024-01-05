package core

type PagedList[T any] struct {
	Size    int `json:"size"`
	From    int `json:"from"`
	Results []T `json:"results"`
}
