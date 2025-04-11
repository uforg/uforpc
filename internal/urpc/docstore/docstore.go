package docstore

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	ErrFileNotFound = os.ErrNotExist
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
	mu        sync.RWMutex             // Protects concurrent access to memCache and diskCache
}

// NewDocstore creates a new Docstore. Read more about at Docstore documentation.
func NewDocstore() *Docstore {
	return &Docstore{
		memCache:  make(map[string]MemCacheFile),
		diskCache: make(map[string]DiskCacheFile),
		mu:        sync.RWMutex{},
	}
}

// normalizePath ensures that the Path is an absolute path and canonicalizes
// it so that it can be used as a key in the cache.
//
// Parameters:
//
//   - relativeToFilePath: An optional absolute file path or URI if the filePath is relative.
//   - filePath: The file path or URI to normalize, if no relativeToFilePath is provided, this should be absolute.
func normalizePath(relativeToFilePath string, filePath string) (string, error) {
	// Convert URI to file path if needed
	filePath = uriToFilePath(filePath)

	if relativeToFilePath != "" {
		// Convert URI to file path if needed
		relativeToFilePath = uriToFilePath(relativeToFilePath)

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

// uriToFilePath converts a URI to a file path.
// It handles both "file://" URIs and regular file paths.
func uriToFilePath(uriOrPath string) string {
	// If it's not a URI, return as is
	if !strings.HasPrefix(uriOrPath, "file://") {
		return filepath.Clean(uriOrPath)
	}

	// Parse the URI
	u, err := url.Parse(uriOrPath)
	if err != nil || u.Scheme != "file" {
		// If parsing fails or it's not a file URI, fall back to simple trimming
		return filepath.Clean(strings.TrimPrefix(uriOrPath, "file://"))
	}

	// Convert the URI path to a file path
	path := u.Path

	// On Windows, handle drive letters correctly
	if len(path) >= 3 && path[0] == '/' && path[2] == ':' {
		// Remove leading slash for Windows drive paths: /C:/path -> C:/path
		path = path[1:]
	}

	return filepath.Clean(path)
}

// OpenInMem opens a file in memory and caches it in the Docstore.
func (d *Docstore) OpenInMem(filePath string, content string) error {
	normFilePath, err := normalizePath("", filePath)
	if err != nil {
		return fmt.Errorf("error normalizing file path: %w", err)
	}

	sum := sha256.Sum256([]byte(content))
	hash := fmt.Sprintf("%x", sum)

	d.mu.Lock()
	defer d.mu.Unlock()

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
	normFilePath, err := normalizePath("", filePath)
	if err != nil {
		return fmt.Errorf("error normalizing file path: %w", err)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.memCache, normFilePath)

	return nil
}

// GetInMemory returns the content of a file in memory, the hash and a boolean
// indicating if the file exists in memory.
//
// Parameters:
//
//   - relativeToFilePath: An optional absolute file path or URI if the filePath is relative.
//   - filePath: The file path or URI to get, if no relativeToFilePath is provided, this should be absolute.
func (d *Docstore) GetInMemory(relativeToFilePath string, filePath string) (string, string, bool, error) {
	normFilePath, err := normalizePath(relativeToFilePath, filePath)
	if err != nil {
		return "", "", false, fmt.Errorf("error normalizing file path: %w", err)
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

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
//   - relativeToFilePath: An optional absolute file path or URI if the filePath is relative.
//   - filePath: The file path or URI to get, if no relativeToFilePath is provided, this should be absolute.
func (d *Docstore) GetFromDisk(relativeToFilePath string, filePath string) (string, string, bool, error) {
	// 1. Normalize the file path
	normFilePath, err := normalizePath(relativeToFilePath, filePath)
	if err != nil {
		return "", "", false, fmt.Errorf("error normalizing file path: %w", err)
	}

	// Use a loop instead of recursion to avoid potential stack overflow
	for {
		// Start with a read lock
		d.mu.RLock()

		// 2. Check if the file exists in diskCache
		cachedFile, ok := d.diskCache[normFilePath]

		if !ok {
			// Not in cache, release read lock
			d.mu.RUnlock()

			// Check if file exists and get info
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

			// Read file content
			content, err := os.ReadFile(normFilePath)
			if err != nil {
				return "", "", false, fmt.Errorf("error reading file: %w", err)
			}

			// Calculate hash
			sum := sha256.Sum256(content)
			hash := fmt.Sprintf("%x", sum)

			// Acquire write lock to update cache
			d.mu.Lock()
			d.diskCache[normFilePath] = DiskCacheFile{
				Content: string(content),
				Hash:    hash,
				Mtime:   fileInfo.ModTime(),
			}
			d.mu.Unlock()

			return string(content), hash, true, nil
		}

		// File is in cache, get file info to check if it's stale
		// We can release the read lock while checking the file
		d.mu.RUnlock()

		fileInfo, err := os.Stat(normFilePath)
		if errors.Is(err, os.ErrNotExist) {
			// File no longer exists, acquire write lock to remove from cache
			d.mu.Lock()
			delete(d.diskCache, normFilePath)
			d.mu.Unlock()
			return "", "", false, nil
		}
		if err != nil {
			return "", "", false, fmt.Errorf("error getting file info: %w", err)
		}
		if fileInfo.IsDir() {
			return "", "", false, fmt.Errorf("file path is a directory: %s", normFilePath)
		}

		// Check if file has been modified
		mtime := fileInfo.ModTime()
		if mtime != cachedFile.Mtime {
			// File has changed, acquire write lock to remove from cache
			d.mu.Lock()
			delete(d.diskCache, normFilePath)
			d.mu.Unlock()
			// Continue the loop to read the updated file
			continue
		}

		// File is in cache and not stale, return cached content
		return cachedFile.Content, cachedFile.Hash, true, nil
	}
}

// GetFileAndHash implements analyzer.FileProvider.GetFileAndHash. It first checks
// try to get the file from memCache using GetInMemory, if not found then try to get
// the file from diskCache using GetFromDisk. If both fails then return an error.
func (d *Docstore) GetFileAndHash(relativeTo string, path string) (string, string, error) {
	content, hash, exists, err := d.GetInMemory(relativeTo, path)
	if err != nil {
		return "", "", fmt.Errorf("error getting file from memCache: %w", err)
	}
	if exists {
		return content, hash, nil
	}

	content, hash, exists, err = d.GetFromDisk(relativeTo, path)
	if err != nil {
		return "", "", fmt.Errorf("error getting file from diskCache: %w", err)
	}
	if exists {
		return content, hash, nil
	}

	// Return a standard error that can be checked with errors.Is
	return "", "", fmt.Errorf("file not found: %s: %w", path, os.ErrNotExist)
}
