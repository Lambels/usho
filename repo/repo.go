package repo

import (
	"context"
	"fmt"
)

var ErrNotFound = fmt.Errorf("couldnt find the specified url")

type Repo interface {
	New(ctx context.Context, in URLRequest) (URLResponse, error)
	Get(ctx context.Context, in string) (URLResponse, error)
	Close() error
}

type URLRequest struct {
	Intial string `json:"initial" validate:"min=7,regexp=[(http(s)?):\/\/(www\.)?a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)"`
}

type URLResponse struct {
	ID      uint64 `json:"id"`
	Initial string `json:"intial"`
	Short   string `json:"short"`
}

type URLKey struct{}
