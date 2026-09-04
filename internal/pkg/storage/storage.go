package storage

import (
	"context"
	"errors"
	"io"
)

var ErrNotFound = errors.New("storage: object not found")

type Object struct {
	Body          io.ReadCloser
	ContentType   string
	ContentLength int64
}

type Getter interface {
	Get(ctx context.Context, key string) (*Object, error)
}
