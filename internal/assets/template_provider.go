package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type FileTemplate struct {
	Content   []byte
	Filename  string
	OutputDir string
}

type ComposePresetTemplates struct {
	Content []byte
	Volume  []string
}

type TemplateProvider struct {
	fileTemplates           map[string]FileTemplate
	composePresetsTemplates map[string]ComposePresetTemplates
}

// all:template will include any hidden file / dir in templates

//go:embed all:templates
var templatesFS embed.FS

// volume configurations for compose presets,
var composeVolumeConfig = map[string][]string{
	"caddy":      {"caddy_data", "caddy_config"},
	"nginx":      nil,
	"postgresql": {"postgresql_data"},
	"redis":      {"redis_data"},
}

const (
	templateEmbedDirectory = "templates"
	composeEmbedDirectory  = "compose"
)

func NewTemplateProvider(generatedDir string) (*TemplateProvider, error) {
	provider := &TemplateProvider{
		fileTemplates:           make(map[string]FileTemplate),
		composePresetsTemplates: make(map[string]ComposePresetTemplates),
	}

	err := provider.loadFileTemplates(generatedDir)
	if err != nil {
		return nil, err
	}
	err = provider.loadComposePresets()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// loading, and mapping the files
func (tp *TemplateProvider) loadFileTemplates(generatedDir string) error {
	// template mappings, id: path
	fileTemplateMappings := map[string]string{
		"caddy":          "files/Caddyfile.tmpl",
		"nginx":          "files/nginx.conf.tmpl",
		"rsync-ignore":   "files/.rsyncignore.tmpl",
		"docker-compose": "files/docker-compose.yaml.tmpl",

		"ansible-setup":     "ansible/setup.yaml.tmpl",
		"ansible-vars":      "ansible/all.yaml.tmpl",
		"ansible-inventory": "ansible/inventory.ini.tmpl",
		"ansible-molecule":  "ansible/converge.yml.tmpl",

		"go-dockerfile": "files/Dockerfile.tmpl",

		"go-actions-ci": "files/go-ci.yaml.tmpl",
	}

	// filename mappings for output id: output_name
	fileNameMappings := map[string]string{
		"caddy":          "Caddyfile",
		"nginx":          "nginx.conf",
		"docker-compose": "docker-compose.yaml",
		"rsync-ignore":   ".rsyncignore",

		"ansible-setup":     "setup.yaml",
		"ansible-vars":      "all.yaml",
		"ansible-inventory": "inventory.ini",
		"ansible-molecule":  "converge.yml",

		"go-dockerfile": "Dockerfile",

		"go-actions-ci": "go-ci.yaml",
	}

	outputDirMappings := map[string]string{
		"caddy":          filepath.Join(generatedDir, "conf"),
		"nginx":          filepath.Join(generatedDir, "conf"),
		"rsync-ignore":   filepath.Join(generatedDir, "conf"),
		"docker-compose": filepath.Join(generatedDir, "conf"),

		"ansible-setup":     filepath.Join(generatedDir, "ansible"),
		"ansible-vars":      filepath.Join(generatedDir, "ansible/group_vars"),
		"ansible-inventory": filepath.Join(generatedDir, "ansible"),
		"ansible-molecule":  filepath.Join(generatedDir, "ansible/molecule/default"),

		"go-dockerfile": filepath.Join(generatedDir, "conf"),

		"go-actions-ci": ".github/workflows",
	}

	subFS, err := fs.Sub(templatesFS, templateEmbedDirectory)
	if err != nil {
		return fmt.Errorf("creating sub filesystem for '%s': %w", templateEmbedDirectory, err)
	}

	for key, path := range fileTemplateMappings {
		content, err := fs.ReadFile(subFS, path)
		if err != nil {
			return fmt.Errorf("failed to read file on %s", path)
		}

		tp.fileTemplates[key] = FileTemplate{
			Content:   content,
			Filename:  fileNameMappings[key],
			OutputDir: outputDirMappings[key],
		}
	}

	return nil
}

func (tp *TemplateProvider) loadComposePresets() error {
	subFS, err := fs.Sub(templatesFS, templateEmbedDirectory)
	if err != nil {
		return fmt.Errorf("creating sub filesystem for '%s': %w", templateEmbedDirectory, err)
	}

	entries, err := fs.ReadDir(subFS, composeEmbedDirectory)
	if err != nil {
		return fmt.Errorf("reading compose directory '%s': %w", composeEmbedDirectory, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tmpl") {
			continue
		}

		// extract preset name from filename (remove .tmpl extension)
		presetName := strings.TrimSuffix(entry.Name(), ".tmpl")

		content, err := fs.ReadFile(subFS, filepath.Join(composeEmbedDirectory, entry.Name()))
		if err != nil {
			return fmt.Errorf("reading compose preset named '%s': %w", entry.Name(), err)
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
