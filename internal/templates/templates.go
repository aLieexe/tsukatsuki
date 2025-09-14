package templates

import (
	_ "embed"
)

// This templates are used for files that directly gives out files
type FileTemplates struct {
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

func (T *FileTemplates) Caddyfile() FileTemplates {
	return FileTemplates{
		Content:  caddyfileContents,
		Filename: "Caddyfile",
	}
}

func (T *FileTemplates) DockerCompose() FileTemplates {
	return FileTemplates{
		Content:  dockerComposeContent,
		Filename: "docker-compose.yaml",
	}
}

func (T *FileTemplates) Dockerfile() FileTemplates {
	return FileTemplates{
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
