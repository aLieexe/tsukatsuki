package templates

import (
	_ "embed"
)

// This templates are used for files that directly gives out files
type FileTemplate struct {
	Content  []byte
	Filename string
}

// *files directory
//

//go:embed files/Caddyfile.tmpl
var caddyfileContents []byte

//go:embed files/docker-compose.yaml.tmpl
var dockerComposeContent []byte

//go:embed files/Dockerfile.tmpl
var dockerfileContent []byte

func (T *FileTemplate) Caddyfile() FileTemplate {
	return FileTemplate{
		Content:  caddyfileContents,
		Filename: "Caddyfile",
	}
}

func (T *FileTemplate) DockerCompose() FileTemplate {
	return FileTemplate{
		Content:  dockerComposeContent,
		Filename: "docker-compose.yaml",
	}
}

func (T *FileTemplate) Dockerfile() FileTemplate {
	return FileTemplate{
		Content:  dockerfileContent,
		Filename: "Dockerfile",
	}
}

// This templates are used for files that is used in FileTemplate
type PresetTemplates struct {
	Content []byte
	Name    string
}

// * compose_presets directory
//

//go:embed compose_presets/caddy.tmpl
var caddyCompose []byte

func (T *PresetTemplates) CaddyCompose() PresetTemplates {
	return PresetTemplates{
		Content: caddyCompose,
		Name:    "caddy",
	}
}
