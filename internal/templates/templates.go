package templates

import (
	_ "embed"
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

//go:embed files/Caddyfile.tmpl
var caddyfileContents []byte

//go:embed files/docker-compose.yaml.tmpl
var dockerComposeContent []byte

//go:embed files/Dockerfile.tmpl
var dockerfileContent []byte

//go:embed files/nginx.conf.tmpl
var nginxConfContent []byte

//go:embed ansible/playbook.yaml.tmpl
var ansiblePlaybookContent []byte

//go:embed ansible/vars.yaml.tmpl
var ansibleVarsContent []byte

//go:embed compose_presets/caddy.tmpl
var caddyCompose []byte

// creates a new template provider with all templates
func NewTemplateProvider() *TemplateProvider {
	provider := &TemplateProvider{
		fileTemplates:           make(map[string]FileTemplate),
		composePresetsTemplates: make(map[string]ComposePresetTemplates),
	}

	// init file templates
	provider.fileTemplates["caddy"] = FileTemplate{
		Content:  caddyfileContents,
		Filename: "Caddyfile",
	}

	provider.fileTemplates["nginx"] = FileTemplate{
		Content:  nginxConfContent,
		Filename: "nginx.conf",
	}

	provider.fileTemplates["dockercompose"] = FileTemplate{
		Content:  dockerComposeContent,
		Filename: "docker-compose.yaml",
	}

	provider.fileTemplates["dockerfile"] = FileTemplate{
		Content:  dockerfileContent,
		Filename: "Dockerfile",
	}

	provider.fileTemplates["ansibleplaybook"] = FileTemplate{
		Content:  ansiblePlaybookContent,
		Filename: "playbook.yaml",
	}

	provider.fileTemplates["ansiblevars"] = FileTemplate{
		Content:  ansibleVarsContent,
		Filename: "vars.yaml",
	}

	// init preset templates
	provider.composePresetsTemplates["caddy"] = ComposePresetTemplates{
		Content: caddyCompose,
		Volume:  []string{"caddy_data", "caddy_config"},
	}

	provider.composePresetsTemplates["nginx"] = ComposePresetTemplates{
		Content: caddyCompose,
		Volume:  nil,
	}

	return provider
}

func (tp *TemplateProvider) GetFileTemplates() map[string]FileTemplate {
	return tp.fileTemplates
}

func (tp *TemplateProvider) GetComposePresetTemplates() map[string]ComposePresetTemplates {
	return tp.composePresetsTemplates
}
