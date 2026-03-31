package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Storage is the interface for object storage operations.
type Storage interface {
	// Upload stores data at the given key and returns the public URL.
	Upload(ctx context.Context, key string, reader io.Reader, contentType string) (url string, err error)
	// Delete removes the object at the given key.
	Delete(ctx context.Context, key string) error
	// PresignedURL generates a time-limited URL for direct access to the object.
	PresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)
}

// LocalStorage implements Storage using the local filesystem.
// Suitable for development and testing environments.
type LocalStorage struct {
	basePath string // directory to store files
	baseURL  string // URL prefix for accessing files
}

// NewLocal creates a LocalStorage that saves files to basePath
// and returns URLs prefixed with baseURL.
func NewLocal(basePath, baseURL string) *LocalStorage {
	return &LocalStorage{basePath: basePath, baseURL: baseURL}
}

// Upload saves the data to a local file and returns the URL.
func (s *LocalStorage) Upload(_ context.Context, key string, reader io.Reader, _ string) (string, error) {
	fullPath := filepath.Join(s.basePath, key)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("storage mkdir: %w", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("storage create: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return "", fmt.Errorf("storage write: %w", err)
	}

	return s.baseURL + "/" + key, nil
}

// Delete removes the file at the given key.
func (s *LocalStorage) Delete(_ context.Context, key string) error {
	fullPath := filepath.Join(s.basePath, key)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("storage delete: %w", err)
	}
	return nil
}

// PresignedURL returns the direct URL for local storage (no expiry enforcement).
func (s *LocalStorage) PresignedURL(_ context.Context, key string, _ time.Duration) (string, error) {
	return s.baseURL + "/" + key, nil
}
