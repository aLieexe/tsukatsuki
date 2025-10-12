package assets

import (
	"embed"
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
			return err
		}

		// calculate relative path from source
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dst, relPath)

		if dirEntry.IsDir() {
			// create directory structure
			return os.MkdirAll(destPath, 0o755)
		}

		return CopyEmbeddedFiles(path, destPath)
	})
}

func CopyEmbeddedFiles(src, dst string) error {
	parentDir := filepath.Dir(dst)
	if err := os.MkdirAll(parentDir, 0o755); err != nil {
		return err
	}

	// read file content from embedded filesystem
	data, err := staticFS.ReadFile(src)
	if err != nil {
		return err
	}

	// write file to destination with appropriate permissions
	return os.WriteFile(dst, data, 0o644)
}
