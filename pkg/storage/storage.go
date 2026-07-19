package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/logger"
)

// Storage defines the file storage interface.
type Storage interface {
	Upload(ctx context.Context, path string, reader io.Reader, contentType string) (string, error)
	Delete(ctx context.Context, path string) error
	GenerateURL(ctx context.Context, path string) (string, error)
}

// FileInfo holds metadata about an uploaded file.
type FileInfo struct {
	Path        string `json:"path"`
	URL         string `json:"url"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}

// New creates a Storage provider based on configuration.
func New(cfg *configs.Config, log *logger.Logger) (Storage, error) {
	switch cfg.Storage.Provider {
	case "local":
		return NewLocalStorage(cfg.Storage.LocalPath, cfg.App.URL, log)
	case "s3":
		return NewS3Storage(&cfg.AWS, log)
	case "r2":
		return NewR2Storage(&cfg.R2, log)
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", cfg.Storage.Provider)
	}
}
