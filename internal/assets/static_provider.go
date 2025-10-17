package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// all:static will include any hidden file / dir in static

//go:embed all:static
var staticFS embed.FS

func CopyEmbeddedDirectory(src, dst string) error {
	return fs.WalkDir(staticFS, src, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking %s: %w", path, err)
		}

		// calculate relative path from source
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("calculating relative path %s: %w", path, err)
		}

		destPath := filepath.Join(dst, relPath)

		if dirEntry.IsDir() {
			// create directory structure
			if err := os.MkdirAll(destPath, 0o755); err != nil {
				return fmt.Errorf("creating directory %s: %w", destPath, err)
			}
			return nil
		}

		return CopyEmbeddedFiles(path, destPath)
	})
}

func CopyEmbeddedFiles(src, dst string) error {
	parentDir := filepath.Dir(dst)
	if err := os.MkdirAll(parentDir, 0o755); err != nil {
		return fmt.Errorf("creating directory %s: %w", parentDir, err)
	}

	// read file content from embedded filesystem
	data, err := staticFS.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading file from embed %s: %w", src, err)
	}

	// write file to destination with appropriate permissions
	return os.WriteFile(dst, data, 0o644)
}
