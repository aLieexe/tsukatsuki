package templates

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"
)

type FileTemplate struct {
	Content  []byte
	Filename string
}

type ComposePresetTemplates struct {
	Content []byte
	Volume  []string
}

type TemplateProvider struct {
	fileTemplates           map[string]FileTemplate
	composePresetsTemplates map[string]ComposePresetTemplates
}

//go:embed files/* ansible/* compose_presets/*
var templatesFS embed.FS

// volume configurations for compose presets
var composeVolumeConfig = map[string][]string{
	"caddy": {"caddy_data", "caddy_config"},
	"nginx": nil,
}

func NewTemplateProvider() (*TemplateProvider, error) {
	provider := &TemplateProvider{
		fileTemplates:           make(map[string]FileTemplate),
		composePresetsTemplates: make(map[string]ComposePresetTemplates),
	}

	err := provider.loadFileTemplates()
	if err != nil {
		return nil, err
	}
	err = provider.loadComposePresets()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (tp *TemplateProvider) loadFileTemplates() error {
	// define template mappings
	templateMappings := map[string]string{
		"caddy":             "files/Caddyfile.tmpl",
		"nginx":             "files/nginx.conf.tmpl",
		"docker-compose":    "files/docker-compose.yaml.tmpl",
		"dockerfile":        "files/Dockerfile.tmpl",
		"ansible-setup":     "ansible/setup.yaml.tmpl",
		"ansible-vars":      "ansible/all.yaml.tmpl",
		"ansible-inventory": "ansible/inventory.ini.tmpl",
	}

	// filename mappings for output
	filenameMappings := map[string]string{
		"caddy":             "Caddyfile",
		"nginx":             "nginx.conf",
		"docker-compose":    "docker-compose.yaml",
		"dockerfile":        "Dockerfile",
		"ansible-setup":     "setup.yaml",
		"ansible-vars":      "all.yaml",
		"ansible-inventory": "inventory.ini",
	}

	for key, path := range templateMappings {
		content, err := templatesFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file on %s", path)
		}

		tp.fileTemplates[key] = FileTemplate{
			Content:  content,
			Filename: filenameMappings[key],
		}
	}

	return nil
}

func (tp *TemplateProvider) loadComposePresets() error {
	// read all files in compose_presets directory
	entries, err := templatesFS.ReadDir("compose_presets")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tmpl") {
			continue
		}

		// extract preset name from filename (remove .tmpl extension)
		presetName := strings.TrimSuffix(entry.Name(), ".tmpl")

		content, err := templatesFS.ReadFile(filepath.Join("compose_presets", entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to read compose preset named %s", entry.Name())
		}

		tp.composePresetsTemplates[presetName] = ComposePresetTemplates{
			Content: content,
			Volume:  composeVolumeConfig[presetName],
		}
	}

	return err
}

func (tp *TemplateProvider) GetFileTemplates() map[string]FileTemplate {
	return tp.fileTemplates
}

func (tp *TemplateProvider) GetComposePresetTemplates() map[string]ComposePresetTemplates {
	return tp.composePresetsTemplates
}
