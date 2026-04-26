package storage

import (
	"context"
	"errors"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

// GCSStore implements ObjectStore using a single Google Cloud Storage bucket.
type GCSStore struct {
	client *storage.Client
	bucket *storage.BucketHandle
}

// NewGCSStore creates a client using Application Default Credentials and
// targets the given bucket name. The caller should call Close when done.
func NewGCSStore(ctx context.Context, bucketName string) (*GCSStore, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("storage: GCS bucket name is empty")
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GCSStore{
		client: client,
		bucket: client.Bucket(bucketName),
	}, nil
}

// Close releases the GCS client. It does not delete bucket data.
func (s *GCSStore) Close() error {
	if s == nil || s.client == nil {
		return nil
	}
	return s.client.Close()
}

// Put implements ObjectStore.
func (s *GCSStore) Put(ctx context.Context, key string, r io.Reader, _ int64, contentType string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if key == "" {
		return ErrInvalidKey
	}

	w := s.bucket.Object(key).NewWriter(ctx)
	if contentType != "" {
		w.ContentType = contentType
	}

	if _, err := io.Copy(w, r); err != nil {
		_ = w.Close()
		return err
	}
	return w.Close()
}

// Get implements ObjectStore.
func (s *GCSStore) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if key == "" {
		return nil, ErrInvalidKey
	}

	rc, err := s.bucket.Object(key).NewReader(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return rc, nil
}

// Delete implements ObjectStore.
func (s *GCSStore) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if key == "" {
		return ErrInvalidKey
	}

	err := s.bucket.Object(key).Delete(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

var _ ObjectStore = (*GCSStore)(nil)
