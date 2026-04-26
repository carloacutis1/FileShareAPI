package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestMemoryStore_PutGetRoundTrip(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStore()
	key := "uploads/doc.txt"
	payload := []byte("hello blob")

	if err := s.Put(ctx, key, bytes.NewReader(payload), int64(len(payload)), "text/plain"); err != nil {
		t.Fatalf("Put: %v", err)
	}

	rc, err := s.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	t.Cleanup(func() { _ = rc.Close() })

	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("content: got %q, want %q", got, payload)
	}
}

func TestMemoryStore_GetNotFound(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStore()

	_, err := s.Get(ctx, "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get missing: got %v, want ErrNotFound", err)
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStore()
	key := "a/b"

	if err := s.Put(ctx, key, strings.NewReader("x"), 1, ""); err != nil {
		t.Fatalf("Put: %v", err)
	}
	if err := s.Delete(ctx, key); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if err := s.Delete(ctx, key); !errors.Is(err, ErrNotFound) {
		t.Fatalf("second Delete: got %v, want ErrNotFound", err)
	}
	_, err := s.Get(ctx, key)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get after delete: got %v, want ErrNotFound", err)
	}
}

func TestMemoryStore_InvalidKey(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStore()

	if err := s.Put(ctx, "", strings.NewReader("x"), 1, ""); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("Put empty key: got %v, want ErrInvalidKey", err)
	}
	if _, err := s.Get(ctx, ""); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("Get empty key: got %v, want ErrInvalidKey", err)
	}
	if err := s.Delete(ctx, ""); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("Delete empty key: got %v, want ErrInvalidKey", err)
	}
}

func TestMemoryStore_CanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	s := NewMemoryStore()

	if err := s.Put(ctx, "k", strings.NewReader("x"), 1, ""); err == nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("Put canceled: got %v, want context.Canceled", err)
	}
	if _, err := s.Get(ctx, "k"); err == nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("Get canceled: got %v, want context.Canceled", err)
	}
	if err := s.Delete(ctx, "k"); err == nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("Delete canceled: got %v, want context.Canceled", err)
	}
}
