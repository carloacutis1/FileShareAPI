package storage

import (
	"context"
	"errors"
	"io"
)

// Common errors for ObjectStore implementations.
var (
	ErrNotFound   = errors.New("storage: object not found")
	ErrInvalidKey = errors.New("storage: key is empty")
)

// ObjectStore abstracts blob storage (e.g. GCS, S3, or an in-memory test double).
// Keys are object paths within a bucket, not full gs:// URLs.
type ObjectStore interface {
	// Put writes or overwrites the object at key. contentType may be empty;
	// implementations that care (e.g. GCS) can set it on the object.
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	// Get returns a ReadCloser for the object at key. The caller must close it.
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	// Delete removes the object at key. ErrNotFound is returned if the key does not exist.
	Delete(ctx context.Context, key string) error
}
