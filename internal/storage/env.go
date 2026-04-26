package storage

import (
	"context"
	"os"
)

// NewBlobStoreFromEnv returns an ObjectStore for the process.
// If GCS_BUCKET is set, it uses Google Cloud Storage (Application Default Credentials).
// Otherwise it returns an in-memory store suitable for local development and tests.
func NewBlobStoreFromEnv(ctx context.Context) (store ObjectStore, cleanup func(), err error) {
	if bucket := os.Getenv("GCS_BUCKET"); bucket != "" {
		gcs, e := NewGCSStore(ctx, bucket)
		if e != nil {
			return nil, nil, e
		}
		return gcs, func() { _ = gcs.Close() }, nil
	}
	return NewMemoryStore(), func() {}, nil
}
