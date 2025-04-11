package docstore

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MemCacheFile is a file that is stored in memory and cached in Docstore.
type MemCacheFile struct {
	Content string
	Hash    string
}

// DiskCacheFile is a file that is stored on disk and cached in memory in Docstore.
type DiskCacheFile struct {
	Content string
	Hash    string
	Mtime   time.Time
}

// Docstore is an implementation of analyzer.FileProvider that caches files in memory.
//
// It has two caches:
//
//   - memCache: Caches files that are opened in memory (can or can't exist on disk).
//   - diskCache: Caches files that are opened on disk (used as fallback of memCache).
type Docstore struct {
	memCache  map[string]MemCacheFile  // Normalized Path -> MemCacheFile
	diskCache map[string]DiskCacheFile // Normalized Path -> DiskCacheFile
}

// NewDocstore creates a new Docstore. Read more about at Docstore documentation.
func NewDocstore() *Docstore {
	return &Docstore{
		memCache:  make(map[string]MemCacheFile),
		diskCache: make(map[string]DiskCacheFile),
	}
}

// normalizePath ensures that the Path is an absolute path and canonicalizes
// it so that it can be used as a key in the cache.
//
// Parameters:
//
//   - relativeToFilePath: An optional absolute file path if the filePath is relative.
//   - filePath: The file path to normalize, if no relativeToFilePath is provided, this should be absolute.
func (d *Docstore) normalizePath(relativeToFilePath string, filePath string) (string, error) {
	filePath = strings.TrimPrefix(filePath, "file://")
	filePath = filepath.Clean(filePath)

	if relativeToFilePath != "" {
		relativeToFilePath = strings.TrimPrefix(relativeToFilePath, "file://")
		relativeToFilePath = filepath.Clean(relativeToFilePath)
		if !filepath.IsAbs(relativeToFilePath) {
			return "", fmt.Errorf("relativeToFilePath must be an absolute path, got %s", relativeToFilePath)
		}

		// Keep only the directory
		relativeToFilePath = filepath.Dir(relativeToFilePath)
	}

	newNormFilePath := filepath.Join(relativeToFilePath, filePath)
	if !filepath.IsAbs(newNormFilePath) {
		return newNormFilePath, fmt.Errorf("file path must be an absolute path, got %s", filePath)
	}

	return newNormFilePath, nil
}

// OpenInMem opens a file in memory and caches it in the Docstore.
func (d *Docstore) OpenInMem(filePath string, content string) error {
	normFilePath, err := d.normalizePath("", filePath)
	if err != nil {
		return fmt.Errorf("error normalizing file path: %w", err)
	}

	sum := sha256.Sum256([]byte(content))
	hash := string(sum[:])

	d.memCache[normFilePath] = MemCacheFile{
		Content: content,
		Hash:    hash,
	}

	// If exists in diskCache then delete it
	// to prioritize the in-memory version
	delete(d.diskCache, normFilePath)

	return nil
}

// ChangeInMem changes the content of a file in memory and caches it in the Docstore.
func (d *Docstore) ChangeInMem(filePath string, content string) error {
	return d.OpenInMem(filePath, content)
}

// CloseInMem closes a file in memory and removes it from the Docstore.
func (d *Docstore) CloseInMem(filePath string) error {
	normFilePath, err := d.normalizePath("", filePath)
	if err != nil {
		return fmt.Errorf("error normalizing file path: %w", err)
	}

	delete(d.memCache, normFilePath)

	return nil
}

// GetInMemory returns the content of a file in memory, the hash and a boolean
// indicating if the file exists in memory.
func (d *Docstore) GetInMemory(filePath string) (string, string, bool, error) {
	normFilePath, err := d.normalizePath("", filePath)
	if err != nil {
		return "", "", false, fmt.Errorf("error normalizing file path: %w", err)
	}

	cachedFile, ok := d.memCache[normFilePath]
	if !ok {
		return "", "", false, nil
	}

	return cachedFile.Content, cachedFile.Hash, true, nil
}

// GetFromDisk returns the content of a file on disk, the hash and a boolean
// indicating if the file exists on disk.
//
// It first checks the diskCache, if found then compares the mtime of the file
// to avoid stale content. If not found in diskCache then it reads the file
// from disk and caches it in diskCache.
//
// Parameters:
//
//   - relativeToFilePath: An optional absolute file path if the filePath is relative.
//   - filePath: The file path to get, if no relativeToFilePath is provided, this should be absolute.
func (d *Docstore) GetFromDisk(relativeToFilePath string, filePath string) (string, string, bool, error) {
	// 1. Normalize the file path
	normFilePath, err := d.normalizePath(relativeToFilePath, filePath)
	if err != nil {
		return "", "", false, fmt.Errorf("error normalizing file path: %w", err)
	}

	// 2. Check if the file exists in diskCache
	cachedFile, ok := d.diskCache[normFilePath]

	// 2.1 If not found in diskCache then read the file from disk and cache it
	if !ok {
		fileInfo, err := os.Stat(normFilePath)
		if errors.Is(err, os.ErrNotExist) {
			return "", "", false, nil
		}
		if err != nil {
			return "", "", false, fmt.Errorf("error getting file info: %w", err)
		}
		if fileInfo.IsDir() {
			return "", "", false, fmt.Errorf("file path is a directory: %s", normFilePath)
		}

		content, err := os.ReadFile(normFilePath)
		if err != nil {
			return "", "", false, fmt.Errorf("error reading file: %w", err)
		}

		sum := sha256.Sum256(content)
		hash := string(sum[:])

		d.diskCache[normFilePath] = DiskCacheFile{
			Content: string(content),
			Hash:    hash,
			Mtime:   fileInfo.ModTime(),
		}

		return string(content), hash, true, nil
	}

	// 2.2 If found in diskCache then compare the mtime of the file to avoid stale content
	fileInfo, err := os.Stat(normFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return "", "", false, nil
	}
	if err != nil {
		return "", "", false, fmt.Errorf("error getting file info: %w", err)
	}
	if fileInfo.IsDir() {
		return "", "", false, fmt.Errorf("file path is a directory: %s", normFilePath)
	}

	// 2.2.1 If mtime is different then invalidate the cache and try again
	mtime := fileInfo.ModTime()
	if mtime != cachedFile.Mtime {
		delete(d.diskCache, normFilePath)
		return d.GetFromDisk(relativeToFilePath, filePath)
	}

	// When cache hit and content not stale, return the cached content
	return cachedFile.Content, cachedFile.Hash, true, nil
}
