package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type StaticFile struct {
	StaticFilePath string
	OutputPath     string
}

type StaticProvider struct {
	StaticFile map[string]StaticFile
}

// all:static will include any hidden file / dir in static
//
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

func NewStaticProvider(generatedDir string) *StaticProvider {
	staticFilePath := map[string]string{
		"ansible-deploy": "static/ansible/deploy.yaml",
		"ansible-config": "static/ansible/ansible.cfg",
	}
	outputPath := map[string]string{
		"ansible-deploy": filepath.Join(generatedDir, "ansible", "deploy.yaml"),
		"ansible-config": filepath.Join(generatedDir, "ansible", "ansible.cfg"),
	}

	provider := &StaticProvider{
		StaticFile: make(map[string]StaticFile),
	}
	for key := range staticFilePath {
		file := StaticFile{
			StaticFilePath: staticFilePath[key],
			OutputPath:     outputPath[key],
		}

		provider.StaticFile[key] = file
	}
	return provider
}
