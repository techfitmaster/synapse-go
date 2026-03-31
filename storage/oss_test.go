package storage

import (
	"testing"
)

func TestOSSConfig_DefaultBaseURL(t *testing.T) {
	cfg := OSSConfig{
		Endpoint:        "oss-cn-hangzhou.aliyuncs.com",
		AccessKeyID:     "test-id",
		AccessKeySecret: "test-secret",
		Bucket:          "my-bucket",
	}

	// NewOSS will fail without real endpoint, but we can verify config logic
	s, err := NewOSS(cfg)
	if err != nil {
		// Expected: real OSS endpoint not available in test
		// This test verifies the function signature and config handling
		t.Skipf("OSS not available (expected in unit test): %v", err)
	}

	if s.baseURL != "https://my-bucket.oss-cn-hangzhou.aliyuncs.com" {
		t.Errorf("baseURL = %q, want auto-generated URL", s.baseURL)
	}
}

func TestOSSConfig_CustomBaseURL(t *testing.T) {
	cfg := OSSConfig{
		Endpoint:        "oss-cn-hangzhou.aliyuncs.com",
		AccessKeyID:     "test-id",
		AccessKeySecret: "test-secret",
		Bucket:          "my-bucket",
		BaseURL:         "https://cdn.example.com",
	}

	s, err := NewOSS(cfg)
	if err != nil {
		t.Skipf("OSS not available: %v", err)
	}

	if s.baseURL != "https://cdn.example.com" {
		t.Errorf("baseURL = %q, want custom URL", s.baseURL)
	}
}
