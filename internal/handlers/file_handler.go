package handlers

import "gitlab.franciscan.edu/spr-26-backend/basement/internal/storage"

// FileHandler serves file HTTP routes and optional blob storage.
type FileHandler struct {
	Blob storage.ObjectStore
}

// NewFileHandler returns a FileHandler that uses the given blob store (must not be nil).
func NewFileHandler(blob storage.ObjectStore) *FileHandler {
	return &FileHandler{Blob: blob}
}
