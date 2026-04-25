package utils

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// CDNConfigInstance holds the loaded CDN configuration
var CDNConfigInstance *CDNConfig

// ServerURL holds the server URL for local uploads
var ServerURL string

// SiteURL holds the frontend site URL
var SiteURL string

// InitCDNConfig initializes CDN config from the main config
func InitCDNConfig(cfg *Config) {
	CDNConfigInstance = &cfg.CDN
	ServerURL = cfg.Server.URL
	SiteURL = cfg.Site.URL
}

// UploadToQiniu uploads a file to Qiniu cloud storage
// fileData: the file content as bytes
// folder: the folder path (e.g., "products", "news")
// ext: file extension including dot (e.g., ".jpg")
// Returns the relative path (e.g., "products/123456_abc.jpg") and error
func UploadToQiniu(fileData []byte, folder, ext string) (string, error) {
	if CDNConfigInstance == nil {
		return "", fmt.Errorf("CDN config not initialized")
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String()[:8], ext)
	
	// Build the key (path in bucket)
	key := fmt.Sprintf("%s/%s", folder, filename)

	// Create mac
	mac := qbox.NewMac(CDNConfigInstance.AccessKey, CDNConfigInstance.SecretKey)

	// Create upload token
	putPolicy := storage.PutPolicy{
		Scope: CDNConfigInstance.Bucket,
	}
	upToken := putPolicy.UploadToken(mac)

	// Create form uploader
	cfg := storage.Config{
		UseHTTPS: true,
	}
	formUploader := storage.NewFormUploader(&cfg)

	// Upload - convert []byte to io.Reader
	err := formUploader.Put(context.Background(), nil, upToken, key, bytes.NewReader(fileData), int64(len(fileData)), nil)
	if err != nil {
		return "", err
	}

	// Return the relative path (without domain)
	return key, nil
}

// GetFullURL returns the full URL for a given path
func GetFullURL(path string) string {
	if path == "" {
		return ""
	}
	
	// If already a full URL, return as-is
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	
	// For old local uploads (/uploads/xxx), prepend server URL
	if strings.HasPrefix(path, "/uploads/") {
		if ServerURL != "" {
			return ServerURL + path
		}
		return path
	}
	
	// For CDN paths, prepend CDN domain
	if CDNConfigInstance == nil {
		return path
	}
	
	// Remove leading slash if present
	path = strings.TrimPrefix(path, "/")
	
	// Ensure domain ends with slash, path does not start with slash
	domain := CDNConfigInstance.Domain
	if !strings.HasSuffix(domain, "/") {
		domain = domain + "/"
	}
	
	return domain + path
}

// GetRelativePath extracts the relative path from a full URL or returns the path as-is
// This is useful for storing paths in the database
func GetRelativePath(urlOrPath string) string {
	if urlOrPath == "" {
		return ""
	}
	
	// If it's not a full URL, return as-is (remove leading slash for consistency)
	if !strings.HasPrefix(urlOrPath, "http://") && !strings.HasPrefix(urlOrPath, "https://") {
		return strings.TrimPrefix(urlOrPath, "/")
	}
	
	// For CDN URLs, extract the path after domain
	if CDNConfigInstance != nil {
		domain := CDNConfigInstance.Domain
		if strings.HasPrefix(urlOrPath, domain+"/") {
			return strings.TrimPrefix(urlOrPath, domain+"/")
		}
	}
	
	// For server URLs, extract the path
	if ServerURL != "" {
		if strings.HasPrefix(urlOrPath, ServerURL) {
			path := strings.TrimPrefix(urlOrPath, ServerURL)
			return strings.TrimPrefix(path, "/")
		}
	}
	
	// If we can't extract, return the original
	return urlOrPath
}
