package filepathutil

import (
	"fmt"
	"path/filepath"
)

// Normalize converts file paths or URIs to absolute, canonical file paths.
// It handles relative paths by resolving them against a base path.
//
// Parameters:
//   - relativeToFilePath: Optional base path for resolving relative paths
//   - filePath: Path to normalize (absolute or relative to relativeToFilePath)
func Normalize(relativeToFilePath string, filePath string) (string, error) {
	// Convert URI to file path if needed
	filePath = FromURI(filePath)

	if relativeToFilePath != "" {
		// Convert URI to file path if needed
		relativeToFilePath = FromURI(relativeToFilePath)

		if !filepath.IsAbs(relativeToFilePath) {
			return "", fmt.Errorf("relativeToFilePath must be an absolute path, got %s", relativeToFilePath)
		}

		// Keep only the directory
		relativeToFilePath = filepath.Dir(relativeToFilePath)
	}

	// Join paths and clean the result
	newNormFilePath := filepath.Clean(filepath.Join(relativeToFilePath, filePath))

	if !filepath.IsAbs(newNormFilePath) {
		return newNormFilePath, fmt.Errorf("file path must be an absolute path, got %s", filePath)
	}

	return newNormFilePath, nil
}
