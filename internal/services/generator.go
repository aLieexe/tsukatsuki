package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/aLieexe/tsukatsuki/internal/templates"
)

// List out all compose services to add in docker-compose.yaml
type ComposeConfig struct {
	Storage  []string
	Services []string
}

// Create output directory, if not exist
// return error if no permission for existing directory
func createOutputDirectory(dir string) error {
	err := os.Mkdir(dir, 0755)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			// if the dir already exists then test the write permission
			testFile := filepath.Join(dir, ".perm_check")
			f, writeErr := os.Create(testFile)
			if writeErr != nil {
				return fmt.Errorf("no write permission in %q: %w", dir, writeErr)
			}
			f.Close()
			os.Remove(testFile) // clean up
			return nil
		}
		// i think parent dir permission also go here? not sure
		return err
	}

	return nil
}

// ? All should just go output to the "tsukatsuki-generated" directory i guess?
func (app *AppConfig) GenerateConfigurationFiles() error {
	outDir := "out"
	err := createOutputDirectory(outDir)
	if err != nil {
		return err
	}

	FileTemplate := templates.FileTemplate{}

	// TODO: Generate the configuration files, and all of its needed configuration (Caddyfile, Dockerfile, etc), Now how tf do i make this able to get scaled, reee
	// Generate Caddyfile test
	tmpl := template.New("tmpl").Option("missingkey=error")
	tmpl = template.Must(tmpl.Parse(string(FileTemplate.Caddyfile().Content)))
	file, err := os.Create(fmt.Sprint(outDir, "/", FileTemplate.Caddyfile().Filename))
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, app)
	if err != nil {
		return err
	}

	// Generate Dockerfile test

	return nil
}

func (app *AppConfig) GenerateCompose() {
	// TODO: Generate docker compose files, and all of its needed configuration (Caddyfile, Dockerfile, etc)

}

// func (app *AppConfig) GenerateAnsibleFiles() {
// 	// TODO: Generate Ansible-playbook, with roles already prepared and configured to the AppConfig

// }
