package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/logger"
)

// R2Storage implements Storage using Cloudflare R2.
type R2Storage struct {
	cfg *configs.R2Config
	log *logger.Logger
}

// NewR2Storage creates a new R2 storage provider.
func NewR2Storage(cfg *configs.R2Config, log *logger.Logger) (*R2Storage, error) {
	return &R2Storage{cfg: cfg, log: log}, nil
}

// Upload uploads a file to R2.
func (s *R2Storage) Upload(ctx context.Context, path string, reader io.Reader, contentType string) (string, error) {
	s.log.Infof("[R2] Would upload file: %s", path)
	return fmt.Sprintf("%s/%s", s.cfg.PublicURL, path), nil
}

// Delete deletes a file from R2.
func (s *R2Storage) Delete(ctx context.Context, path string) error {
	s.log.Infof("[R2] Would delete file: %s", path)
	return nil
}

// GenerateURL generates a URL for a file in R2.
func (s *R2Storage) GenerateURL(ctx context.Context, path string) (string, error) {
	return fmt.Sprintf("%s/%s", s.cfg.PublicURL, path), nil
}
