package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSConfig holds configuration for Aliyun OSS.
type OSSConfig struct {
	Endpoint        string // e.g. "oss-cn-hangzhou.aliyuncs.com"
	AccessKeyID     string
	AccessKeySecret string
	Bucket          string
	BaseURL         string // public URL prefix, e.g. "https://bucket.oss-cn-hangzhou.aliyuncs.com"
}

// OSSStorage implements Storage using Aliyun OSS.
type OSSStorage struct {
	bucket  *oss.Bucket
	baseURL string
}

// NewOSS creates a Storage backed by Aliyun OSS.
func NewOSS(cfg OSSConfig) (*OSSStorage, error) {
	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("storage oss client: %w", err)
	}

	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("storage oss bucket: %w", err)
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = fmt.Sprintf("https://%s.%s", cfg.Bucket, cfg.Endpoint)
	}

	return &OSSStorage{bucket: bucket, baseURL: baseURL}, nil
}

// Upload stores data in OSS and returns the public URL.
func (s *OSSStorage) Upload(_ context.Context, key string, reader io.Reader, contentType string) (string, error) {
	options := []oss.Option{}
	if contentType != "" {
		options = append(options, oss.ContentType(contentType))
	}

	if err := s.bucket.PutObject(key, reader, options...); err != nil {
		return "", fmt.Errorf("storage oss upload: %w", err)
	}

	return s.baseURL + "/" + key, nil
}

// Delete removes an object from OSS.
func (s *OSSStorage) Delete(_ context.Context, key string) error {
	if err := s.bucket.DeleteObject(key); err != nil {
		return fmt.Errorf("storage oss delete: %w", err)
	}
	return nil
}

// PresignedURL generates a time-limited signed URL for direct access.
func (s *OSSStorage) PresignedURL(_ context.Context, key string, expiry time.Duration) (string, error) {
	url, err := s.bucket.SignURL(key, oss.HTTPGet, int64(expiry.Seconds()))
	if err != nil {
		return "", fmt.Errorf("storage oss presign: %w", err)
	}
	return url, nil
}
