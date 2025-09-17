package templates

import (
	_ "embed"
)

type FileTemplate struct {
	Content  []byte
	Filename string
}

type PresetTemplates struct {
	Content []byte
	Name    string
}

type TemplateProvider struct {
	fileTemplates   map[string]FileTemplate
	presetTemplates map[string]PresetTemplates
}

//go:embed files/Caddyfile.tmpl
var caddyfileContents []byte

//go:embed files/docker-compose.yaml.tmpl
var dockerComposeContent []byte

//go:embed files/Dockerfile.tmpl
var dockerfileContent []byte

//go:embed files/nginx.conf.tmpl
var nginxConfContent []byte

//go:embed compose_presets/caddy.tmpl
var caddyCompose []byte

// creates a new template provider with all templates
func NewTemplateProvider() *TemplateProvider {
	provider := &TemplateProvider{
		fileTemplates:   make(map[string]FileTemplate),
		presetTemplates: make(map[string]PresetTemplates),
	}

	// init file templates
	provider.fileTemplates["caddy"] = FileTemplate{
		Content:  caddyfileContents,
		Filename: "Caddyfile",
	}

	provider.fileTemplates["dockerCompose"] = FileTemplate{
		Content:  dockerComposeContent,
		Filename: "docker-compose.yaml",
	}

	provider.fileTemplates["dockerfile"] = FileTemplate{
		Content:  dockerfileContent,
		Filename: "Dockerfile",
	}

	provider.fileTemplates["nginx"] = FileTemplate{
		Content:  nginxConfContent,
		Filename: "nginx.conf",
	}

	// init preset templates
	provider.presetTemplates["caddyCompose"] = PresetTemplates{
		Content: caddyCompose,
		Name:    "caddy",
	}

	return provider
}

func (tp *TemplateProvider) GetFileTemplates() map[string]FileTemplate {
	return tp.fileTemplates
}

func (tp *TemplateProvider) GetPresetTemplates() map[string]PresetTemplates {
	return tp.presetTemplates
}

func (tp *TemplateProvider) GetFileTemplate(name string) (FileTemplate, bool) {
	template, exists := tp.fileTemplates[name]
	return template, exists
}

func (tp *TemplateProvider) GetPresetTemplate(name string) (PresetTemplates, bool) {
	template, exists := tp.presetTemplates[name]
	return template, exists
}
