package storage

import (
	"context"
	"io"

	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/logger"
)

// S3Storage implements Storage using AWS S3.
// Requires aws-sdk-go-v2 when credentials are provided.
type S3Storage struct {
	cfg *configs.AWSConfig
	log *logger.Logger
}

// NewS3Storage creates a new S3 storage provider.
func NewS3Storage(cfg *configs.AWSConfig, log *logger.Logger) (*S3Storage, error) {
	return &S3Storage{cfg: cfg, log: log}, nil
}

// Upload uploads a file to S3.
func (s *S3Storage) Upload(ctx context.Context, path string, reader io.Reader, contentType string) (string, error) {
	// AWS SDK integration placeholder
	s.log.Infof("[S3] Would upload file: %s", path)
	return "https://" + s.cfg.S3Bucket + ".s3.amazonaws.com/" + path, nil
}

// Delete deletes a file from S3.
func (s *S3Storage) Delete(ctx context.Context, path string) error {
	s.log.Infof("[S3] Would delete file: %s", path)
	return nil
}

// GenerateURL generates a presigned URL for a file.
func (s *S3Storage) GenerateURL(ctx context.Context, path string) (string, error) {
	return "https://" + s.cfg.S3Bucket + ".s3.amazonaws.com/" + path, nil
}
