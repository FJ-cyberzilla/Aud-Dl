package gate

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidURL       = errors.New("validator: input query or URL is invalid")
	ErrUnsafePath       = errors.New("validator: directory traversal attempt detected")
	ErrInsufficientDisk = errors.New("validator: insufficient target vault storage space")
	ErrEmptyQuery       = errors.New("validator: search query cannot be empty")
)

type HubValidator struct {
	MaxFilenameLength int
	MinFreeDiskBytes  uint64
}

func NewHubValidator() *HubValidator {
	return &HubValidator{
		MaxFilenameLength: 255,
		MinFreeDiskBytes:  500 * 1024 * 1024, // 500 MB required buffer
	}
}

// ValidateQuery inspects raw user search queries or links prior to processing
func (hv *HubValidator) ValidateQuery(query string) error {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return ErrEmptyQuery
	}

	// Check if input is a URL and enforce HTTPS scheme
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		parsedURL, err := url.Parse(trimmed)
		if err != nil || parsedURL.Host == "" {
			return fmt.Errorf("%w: failed to parse endpoint", ErrInvalidURL)
		}
		if parsedURL.Scheme != "https" {
			return fmt.Errorf("%w: insecure HTTP protocol rejected", ErrInvalidURL)
		}
	}

	return nil
}

// ValidateSanitizeFilename prevents path traversal attacks (e.g., ../../) and strips illegal characters
func (hv *HubValidator) ValidateSanitizeFilename(filename string) (string, error) {
	clean := filepath.Base(filepath.Clean(filename))
	
	// Prevent directory traversal escape
	if strings.Contains(clean, "..") || strings.ContainsAny(clean, `\/:*?"<>|`) {
		// Clean illegal characters
		clean = strings.Map(func(r rune) rune {
			if strings.ContainsRune(`\/:*?"<>|`, r) {
				return '_'
			}
			return r
		}, clean)
	}

	if len(clean) > hv.MaxFilenameLength {
		clean = clean[:hv.MaxFilenameLength]
	}

	if clean == "" || clean == "." {
		return "", ErrUnsafePath
	}

	return clean, nil
}

// ValidateVaultTarget verifies target path existence and write permissions
func (hv *HubValidator) ValidateVaultTarget(targetPath string) error {
	info, err := os.Stat(targetPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("validator: target path '%s' does not exist", targetPath)
	}
	if !info.IsDir() {
		return fmt.Errorf("validator: target path '%s' is not a directory", targetPath)
	}

	return nil
}
