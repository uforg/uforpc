package filepathutil

import (
	"fmt"
	"path/filepath"
)

// ToAbsolutePath converts a relative path to an absolute path
// using the current working directory as the base.
//
// If relativePath is already absolute, it is returned as is.
func ToAbsolutePath(relativePath string) (string, error) {
	// Get working dir
	wd, err := filepath.Abs(".")
	if err != nil {
		return relativePath, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Join and normalize the working dir and the relative path
	return Normalize(wd, relativePath)
}
