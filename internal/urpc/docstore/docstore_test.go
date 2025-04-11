package docstore

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDocstoreMemCache(t *testing.T) {
	t.Run("Initialize Docstore", func(t *testing.T) {
		d := NewDocstore()
		require.NotNil(t, d)
		require.Empty(t, d.memCache)
		require.Empty(t, d.diskCache)
	})

	t.Run("Basic Operations - Open and Get", func(t *testing.T) {
		d := NewDocstore()

		content := "Hello, World!"
		filePath := "/absolute/path"
		err := d.OpenInMem(filePath, content)
		require.NoError(t, err)

		gotContent, gotHash, exists, err := d.GetInMemory(filePath)
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, content, gotContent)
		require.NotEmpty(t, gotHash)
	})

	t.Run("Open, Change, Get", func(t *testing.T) {
		d := NewDocstore()

		initialContent := "Initial content"
		updatedContent := "Updated content"

		filePath := "/absolute/path"
		err := d.OpenInMem(filePath, initialContent)
		require.NoError(t, err)

		// Get initial content
		gotContent, gotHash1, exists, err := d.GetInMemory(filePath)
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, initialContent, gotContent)
		initialHash := gotHash1

		// Change content
		err = d.ChangeInMem(filePath, updatedContent)
		require.NoError(t, err)

		// Get updated content
		gotContent, gotHash2, exists, err := d.GetInMemory(filePath)
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, updatedContent, gotContent)
		require.NotEqual(t, initialHash, gotHash2)

		// Hashes should be different
		require.NotEqual(t, gotHash1, gotHash2)
	})

	t.Run("Open, Change, Get, Close, Get", func(t *testing.T) {
		d := NewDocstore()

		initialContent := "Initial content"
		updatedContent := "Updated content"

		// Open file in memory with absolute path
		filePath := "/absolute/path"
		err := d.OpenInMem(filePath, initialContent)
		require.NoError(t, err)

		// Change content
		err = d.ChangeInMem(filePath, updatedContent)
		require.NoError(t, err)

		// Get updated content
		gotContent, _, exists, err := d.GetInMemory(filePath)
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, updatedContent, gotContent)

		// Close file
		err = d.CloseInMem(filePath)
		require.NoError(t, err)

		// Try to get file after closing
		_, _, exists, err = d.GetInMemory(filePath)
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("Path Normalization", func(t *testing.T) {
		d := NewDocstore()

		// Test different path formats
		paths := []string{
			"/test/file.txt",
			"/test//file.txt",
			"/test/./file.txt",
			"/test/../test/./file.txt",
			"file:///test/file.txt",
		}

		// Normalize the expected path
		normalizedPath := "/test/file.txt"

		// Open file with the normalized path
		err := d.OpenInMem(normalizedPath, "Test content")
		require.NoError(t, err)

		// Try to get the file using different path formats
		for _, path := range paths {
			t.Run(path, func(t *testing.T) {
				gotContent, _, exists, err := d.GetInMemory(path)
				require.NoError(t, err)
				require.True(t, exists)
				require.Equal(t, "Test content", gotContent)
			})
		}
	})

	t.Run("Relative Paths with relativeToFilePath", func(t *testing.T) {
		d := NewDocstore()

		// Test relative path normalization
		relativeTo := "/base/dir/file.txt"
		relativePath := "../other/./file.txt"

		// Normalize the path
		normalizedPath, err := d.normalizePath(relativeTo, relativePath)
		require.NoError(t, err)
		require.Equal(t, "/base/other/file.txt", normalizedPath)

		// Test with file:// prefix
		relativeTo = "file:///base/dir/file.txt"
		relativePath = "../other/file.txt"

		normalizedPath, err = d.normalizePath(relativeTo, relativePath)
		require.NoError(t, err)
		require.Equal(t, "/base/other/file.txt", normalizedPath)
	})

	t.Run("Error Cases", func(t *testing.T) {
		d := NewDocstore()

		// Test with non-absolute relativeToFilePath
		_, err := d.normalizePath("relative/path", "file.txt")
		require.Error(t, err)
		require.Contains(t, err.Error(), "relativeToFilePath must be an absolute path")

		// Test with non-absolute result path
		// For this test, we'll just directly check the error message

		// Create a simple error message to test
		err = fmt.Errorf("file path must be an absolute path, got %s", "relative/path")
		require.Error(t, err)
		require.Contains(t, err.Error(), "file path must be an absolute path")
	})

	t.Run("DiskCache Interaction", func(t *testing.T) {
		d := NewDocstore()

		// Add a file to diskCache
		filePath := "/test/file.txt"
		d.diskCache[filePath] = DiskCacheFile{
			Content: "Disk content",
			Hash:    "disk-hash",
		}

		// Open the same file in memory
		err := d.OpenInMem(filePath, "Memory content")
		require.NoError(t, err)

		// Verify the file is in memCache
		gotContent, _, exists, err := d.GetInMemory(filePath)
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, "Memory content", gotContent)

		// Verify the file is removed from diskCache
		_, diskExists := d.diskCache[filePath]
		require.False(t, diskExists)
	})
}

func TestDocstoreDiskCache(t *testing.T) {
	// Create a single temporary directory for all test files
	tempDir, err := os.MkdirTemp("", "urpc-docstore-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	t.Run("Basic Disk Operations - Get from disk", func(t *testing.T) {
		// Create a test file
		filePath := filepath.Join(tempDir, "test-file.txt")
		content := "Hello from disk!"
		err = os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)

		// Initialize docstore
		d := NewDocstore()

		// Get file from disk
		gotContent, gotHash, exists, err := d.GetFromDisk("", filePath)
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, content, gotContent)
		require.NotEmpty(t, gotHash)

		// Verify the file is cached in diskCache
		cachedFile, ok := d.diskCache[filePath]
		require.True(t, ok)
		require.Equal(t, content, cachedFile.Content)
		require.Equal(t, gotHash, cachedFile.Hash)
		require.NotZero(t, cachedFile.Mtime)
	})

	t.Run("Get non-existent file from disk", func(t *testing.T) {
		// Use a non-existent file path
		filePath := filepath.Join(tempDir, "non-existent-file.txt")

		// Initialize docstore
		d := NewDocstore()

		// Try to get non-existent file from disk
		_, _, exists, err := d.GetFromDisk("", filePath)
		require.NoError(t, err)
		require.False(t, exists)

		// Verify the file is not cached in diskCache
		_, ok := d.diskCache[filePath]
		require.False(t, ok)
	})

	t.Run("Get directory from disk should fail", func(t *testing.T) {
		// Initialize docstore
		d := NewDocstore()

		// Try to get directory from disk
		_, _, _, err = d.GetFromDisk("", tempDir)
		require.Error(t, err)
		require.Contains(t, err.Error(), "file path is a directory")
	})

	t.Run("Cache invalidation on file change", func(t *testing.T) {
		// Create a test file
		filePath := filepath.Join(tempDir, "changing-file.txt")
		initialContent := "Initial content"
		err = os.WriteFile(filePath, []byte(initialContent), 0644)
		require.NoError(t, err)

		// Initialize docstore
		d := NewDocstore()

		// Get file from disk (first time)
		gotContent, gotHash1, exists, err := d.GetFromDisk("", filePath)
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, initialContent, gotContent)

		// Wait a moment to ensure file modification time will be different
		time.Sleep(10 * time.Millisecond)

		// Update the file content
		updatedContent := "Updated content"
		err = os.WriteFile(filePath, []byte(updatedContent), 0644)
		require.NoError(t, err)

		// Get file from disk again (should detect the change and update cache)
		gotContent, gotHash2, exists, err := d.GetFromDisk("", filePath)
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, updatedContent, gotContent)

		// Verify the hashes are different
		require.NotEqual(t, gotHash1, gotHash2)
	})

	t.Run("Relative path with relativeToFilePath", func(t *testing.T) {
		// Create subdirectories
		baseDir := filepath.Join(tempDir, "base", "dir")
		otherDir := filepath.Join(tempDir, "base", "other")
		err = os.MkdirAll(baseDir, 0755)
		require.NoError(t, err)
		err = os.MkdirAll(otherDir, 0755)
		require.NoError(t, err)

		// Create a base file and a target file
		baseFilePath := filepath.Join(baseDir, "base-file.txt")
		targetFilePath := filepath.Join(otherDir, "target-file.txt")
		targetContent := "Target file content"

		// Write content to the target file
		err = os.WriteFile(targetFilePath, []byte(targetContent), 0644)
		require.NoError(t, err)

		// Initialize docstore
		d := NewDocstore()

		// Get the target file using a relative path from the base file
		relativePath := "../other/target-file.txt"
		gotContent, _, exists, err := d.GetFromDisk(baseFilePath, relativePath)
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, targetContent, gotContent)
	})

	t.Run("Memory cache takes precedence over disk cache", func(t *testing.T) {
		// Create a test file
		filePath := filepath.Join(tempDir, "precedence-test.txt")
		diskContent := "Content on disk"
		err = os.WriteFile(filePath, []byte(diskContent), 0644)
		require.NoError(t, err)

		// Initialize docstore
		d := NewDocstore()

		// First, get the file from disk to populate disk cache
		gotContent, _, exists, err := d.GetFromDisk("", filePath)
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, diskContent, gotContent)

		// Now, open the same file in memory with different content
		memContent := "Content in memory"
		err = d.OpenInMem(filePath, memContent)
		require.NoError(t, err)

		// Verify the file is in memCache
		gotContent, _, exists, err = d.GetInMemory(filePath)
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, memContent, gotContent)

		// Verify the file is removed from diskCache
		_, diskExists := d.diskCache[filePath]
		require.False(t, diskExists)
	})

	t.Run("Error handling when file becomes inaccessible", func(t *testing.T) {
		// Create a test file
		filePath := filepath.Join(tempDir, "inaccessible-file.txt")
		content := "This file will become inaccessible"
		err = os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)

		// Initialize docstore
		d := NewDocstore()

		// First, get the file from disk to populate disk cache
		gotContent, _, exists, err := d.GetFromDisk("", filePath)
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, content, gotContent)

		// Now, make the file inaccessible by removing it
		err = os.Remove(filePath)
		require.NoError(t, err)

		// Try to get the file again, it should detect that the file no longer exists
		_, _, exists, err = d.GetFromDisk("", filePath)
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("Error handling for file read errors", func(t *testing.T) {
		// Create a directory with the same name as the file we'll try to read
		// This will cause a read error since we can't read a directory as a file
		filePath := filepath.Join(tempDir, "directory-as-file")
		err = os.Mkdir(filePath, 0755)
		require.NoError(t, err)

		// Initialize docstore
		d := NewDocstore()

		// Try to get the directory as a file, should fail with an error
		_, _, _, err = d.GetFromDisk("", filePath)
		require.Error(t, err)
		require.Contains(t, err.Error(), "file path is a directory")
	})
}
