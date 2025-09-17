package services

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

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

// This is only for standard file, i think ansible files can also be here? Not sure, but most likely yes
func (app *AppConfig) generateStandardTemplate(fileTemplate *templates.FileTemplate, templateName, outDir string) error {
	tmpl, err := template.New(templateName).Option("missingkey=error").Parse(string(fileTemplate.Content))
	if err != nil {
		return fmt.Errorf("error parsing template %s: %w", templateName, err)
	}

	// create output file
	filePath := filepath.Join(outDir, fileTemplate.Filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", filePath, err)
	}
	defer file.Close()

	// execute template with app context
	if err := tmpl.Execute(file, app); err != nil {
		return fmt.Errorf("error executing template %s: %w", templateName, err)
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

	templateProvider := templates.NewTemplateProvider()

	templateNeeded := []string{"dockerfile", "caddy", "nginx"}

	for _, templateName := range templateNeeded {
		fileTemplate := templateProvider.GetFileTemplates()[templateName]
		if err := app.generateStandardTemplate(&fileTemplate, templateName, outDir); err != nil {
			return err
		}
	}

	return nil
}

// TODO: Compose is a little special, so maybe later?
func (app *AppConfig) GenerateCompose() {
}
