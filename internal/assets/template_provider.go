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
	EnvVar  []EnvVar
}

type TemplateProvider struct {
	fileTemplates           map[string]FileTemplate
	composePresetsTemplates map[string]ComposePresetTemplates
}

type EnvVar struct {
	Name    string
	Default string
}

// all:template will include any hidden file / dir in templates

//go:embed all:templates
var templatesFS embed.FS

// volume configurations for compose presets,
var composeVolumeConfig = map[string][]string{
	"caddy":         {"caddy_data", "caddy_config"},
	"caddy-proxy":   {"caddy_data", "caddy_config"},
	"traefik-proxy": {"traefik_data"},
	"postgresql":    {"postgresql_data"},
	"redis":         {"redis_data"},
	"minio":         {"minio_data"},
	"rabbitmq":      {"rabbitmq_data", "rabbitmq_log"},
}

var servicesEnvVar = map[string][]EnvVar{
	"postgresql": {
		EnvVar{Name: "POSTGRES_USER", Default: "user"},
		EnvVar{Name: "POSTGRES_PASSWORD", Default: "12345678"},
		EnvVar{Name: "POSTGRES_DB", Default: "database"},
	},
	"redis": nil,
	"caddy": nil,
	"minio": {
		EnvVar{Name: "MINIO_ROOT_USER", Default: "root"},
		EnvVar{Name: "MINIO_ROOT_PASSWORD", Default: "12345678"},
	},
	"rabbitmq": {
		EnvVar{Name: "RABBITMQ_DEFAULT_USER", Default: "user"},
		EnvVar{Name: "RABBITMQ_DEFAULT_PASS", Default: "12345678"},
		EnvVar{Name: "RABBITMQ_DEFAULT_VHOST", Default: "vhost"},
	},
}

const (
	templateEmbedDirectory = "templates"
	composeEmbedDirectory  = "compose"
)

func NewTemplateProvider(generatedDir string, projectName string) (*TemplateProvider, error) {
	provider := &TemplateProvider{
		fileTemplates:           make(map[string]FileTemplate),
		composePresetsTemplates: make(map[string]ComposePresetTemplates),
	}

	err := provider.loadFileTemplates(generatedDir, projectName)
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
func (tp *TemplateProvider) loadFileTemplates(generatedDir string, projectName string) error {
	// template mappings, id: path
	fileTemplateMappings := map[string]string{
		"proxy-caddy-entry":   "files/Caddyfile.tmpl",
		"proxy-caddy-project": "files/project.caddy.tmpl",

		"rsync-ignore": "files/.rsyncignore.tmpl",

		"compose-project": "files/project-compose.yaml.tmpl",
		"compose-proxy":   "files/proxy-compose.yaml.tmpl",

		"ansible-setup":     "ansible/setup.yaml.tmpl",
		"ansible-vars":      "ansible/all.yaml.tmpl",
		"ansible-inventory": "ansible/inventory.ini.tmpl",
		"ansible-proxy":     "ansible/proxy.main.yaml.tmpl",

		"dockerfile-go":   "files/go-dockerfile.tmpl",
		"dockerfile-node": "files/node-dockerfile.tmpl",

		"actions-ci-go":   "files/go-ci.yaml.tmpl",
		"actions-ci-node": "files/node-ci.yaml.tmpl",

		"actions-cd-docker": "files/docker-cd.yaml.tmpl",
	}

	// filename mappings for output id: output_name
	fileNameMappings := map[string]string{
		"proxy-caddy-entry":   "Caddyfile",
		"proxy-caddy-project": fmt.Sprintf("%s.caddy", projectName),

		"rsync-ignore": ".rsyncignore",

		"compose-project": "project-compose.yaml",
		"compose-proxy":   "proxy-compose.yaml",

		"ansible-setup":     "setup.yaml",
		"ansible-vars":      "all.yaml",
		"ansible-inventory": "inventory.ini",
		"ansible-proxy":     "main.yaml",

		"dockerfile-go":   "Dockerfile",
		"dockerfile-node": "Dockerfile",

		"actions-ci-go":   "CI.yaml",
		"actions-ci-node": "CI.yaml",

		"actions-cd-docker": "CD.yaml",
	}

	outputDirMappings := map[string]string{
		"proxy-caddy-entry":   filepath.Join(generatedDir, "proxy"),
		"proxy-caddy-project": filepath.Join(generatedDir, "proxy/sites"),

		"rsync-ignore": filepath.Join(generatedDir, "conf"),

		"compose-proxy":   filepath.Join(generatedDir, "proxy"),
		"compose-project": filepath.Join(generatedDir, "conf"),

		"ansible-setup":     filepath.Join(generatedDir, "ansible"),
		"ansible-vars":      filepath.Join(generatedDir, "ansible/group_vars"),
		"ansible-inventory": filepath.Join(generatedDir, "ansible"),
		"ansible-proxy":     filepath.Join(generatedDir, "ansible", "roles", "proxy", "tasks"),

		"dockerfile-go":   filepath.Join(generatedDir, "conf"),
		"dockerfile-node": filepath.Join(generatedDir, "conf"),

		"actions-ci-go":   ".github/workflows",
		"actions-ci-node": ".github/workflows",

		"actions-cd-docker": ".github/workflows",
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
			EnvVar:  servicesEnvVar[presetName],
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
