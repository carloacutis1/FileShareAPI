package storage

import (
	"context"
	"testing"
)

func TestNewGCSStore_emptyBucketName(t *testing.T) {
	_, err := NewGCSStore(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty bucket name")
	}
}
