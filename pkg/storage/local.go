package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/web3sphere/backend/pkg/logger"
)

// LocalStorage implements Storage using the local filesystem.
type LocalStorage struct {
	basePath string
	baseURL  string
	log      *logger.Logger
}

// NewLocalStorage creates a new local storage provider.
func NewLocalStorage(basePath, baseURL string, log *logger.Logger) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &LocalStorage{basePath: basePath, baseURL: baseURL, log: log}, nil
}

// Upload saves a file to the local filesystem.
func (s *LocalStorage) Upload(ctx context.Context, path string, reader io.Reader, contentType string) (string, error) {
	fullPath := filepath.Join(s.basePath, path)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	url := fmt.Sprintf("%s/uploads/%s", s.baseURL, path)
	s.log.Infof("File uploaded: %s", fullPath)
	return url, nil
}

// Delete removes a file from the local filesystem.
func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	s.log.Infof("File deleted: %s", fullPath)
	return nil
}

// GenerateURL returns the URL for a stored file.
func (s *LocalStorage) GenerateURL(ctx context.Context, path string) (string, error) {
	return fmt.Sprintf("%s/uploads/%s", s.baseURL, path), nil
}
