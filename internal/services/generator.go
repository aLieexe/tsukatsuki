package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/aLieexe/tsukatsuki/internal/assets"
)

// List out all compose services to add in docker-compose.yaml
type ComposeConfig struct {
	Storage  []string
	Services []string
}

func (app *AppConfig) GenerateDeploymentFiles() error {
	res := make([]string, 0)
	res = append(res, app.Webserver)

	err := app.GenerateCompose(res, filepath.Join(app.OutputDir, "conf"))
	if err != nil {
		return err
	}

	err = app.GenerateAnsibleFiles(res, filepath.Join(app.OutputDir, "ansible"))
	if err != nil {
		return err
	}

	res = append(res, "dockerfile")
	res = append(res, "rsync-ignore")

	err = app.GenerateConfigurationFiles(res, filepath.Join(app.OutputDir, "conf"))
	if err != nil {
		return err
	}

	return nil
}

func (app *AppConfig) GenerateAnsibleFiles(serviceList []string, outDir string) error {
	err := createOutputDirectory(outDir)
	if err != nil {
		return err
	}

	playbookData := struct {
		Roles []string
	}{
		Roles: make([]string, 0),
	}

	playbookData.Roles = append(playbookData.Roles, "common")
	playbookData.Roles = append(playbookData.Roles, "docker")
	playbookData.Roles = append(playbookData.Roles, serviceList...)

	templateProvider, err := assets.NewTemplateProvider()
	if err != nil {
		return err
	}

	fileTemplate := templateProvider.GetFileTemplates()["ansible-setup"]
	if err := generateStandardTemplate(&fileTemplate, "setup-playbook", outDir, playbookData); err != nil {
		return err
	}

	fileTemplate = templateProvider.GetFileTemplates()["ansible-inventory"]
	if err := generateStandardTemplate(&fileTemplate, "inventory", outDir, app); err != nil {
		return err
	}

	varsDir := filepath.Join(outDir, "/group_vars")
	err = createOutputDirectory(varsDir)
	if err != nil {
		return err
	}

	fileTemplate = templateProvider.GetFileTemplates()["ansible-vars"]
	if err := generateStandardTemplate(&fileTemplate, "ansible-vars", varsDir, app); err != nil {
		return err
	}

	// Copy the file that we wont need to template
	err = copyFile("static/ansible/ansible.cfg", filepath.Join(outDir, "ansible.cfg"))
	if err != nil {
		return err
	}

	err = copyFile("static/ansible/deploy.yaml", filepath.Join(outDir, "deploy.yaml"))
	if err != nil {
		return err
	}

	rolesSrcDir := "static/ansible/roles"
	rolesDstDir := filepath.Join(outDir, "/roles")

	playbookData.Roles = append(playbookData.Roles, "deployment")

	for _, role := range playbookData.Roles {
		src := filepath.Join(rolesSrcDir, role)
		dst := filepath.Join(rolesDstDir, role)

		if err := copyDir(src, dst); err != nil {
			return err
		}
	}

	return nil
}

// ? All should just go output to the "tsukatsuki-generated" directory i guess?
func (app *AppConfig) GenerateConfigurationFiles(templateNeeded []string, outDir string) error {
	err := createOutputDirectory(outDir)
	if err != nil {
		return err
	}

	templateProvider, err := assets.NewTemplateProvider()
	if err != nil {
		return err
	}

	for _, templateName := range templateNeeded {
		fileTemplate := templateProvider.GetFileTemplates()[templateName]
		if err := generateStandardTemplate(&fileTemplate, templateName, outDir, app); err != nil {
			return err
		}
	}
	return nil
}

// TODO: ADD MORE PRESETS, TEST IT
func (app *AppConfig) GenerateCompose(presetNeeded []string, outDir string) error {
	err := createOutputDirectory(outDir)
	if err != nil {
		return err
	}

	templateProvider, err := assets.NewTemplateProvider()
	if err != nil {
		return err
	}
	composeTemplate := templateProvider.GetFileTemplates()["docker-compose"]

	tmpl, err := template.New("docker-compose").Option("missingkey=error").Parse(string(composeTemplate.Content))
	if err != nil {
		return fmt.Errorf("error parsing template %s: %w", composeTemplate.Filename, err)
	}

	// create output file
	filePath := filepath.Join(outDir, composeTemplate.Filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", filePath, err)
	}

	defer func() {
		if closeError := file.Close(); closeError != nil {
			if err == nil {
				err = closeError
			}
		}
	}()

	// temp to combine all the presets,
	templateData := struct {
		Service []string
		Volumes []string
	}{
		Service: []string{},
		Volumes: []string{},
	}

	presetProvider := templateProvider.GetComposePresetTemplates()
	for _, presetName := range presetNeeded {
		if preset, exists := presetProvider[presetName]; exists {
			// add the services
			serviceDefinition := string(preset.Content)
			templateData.Service = append(templateData.Service, serviceDefinition)

			// add volumes from preset
			if preset.Volume != nil {
				templateData.Volumes = append(templateData.Volumes, preset.Volume...)
			}
		}
	}

	err = tmpl.Execute(file, templateData)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

// Create output directory, if not exist
// return error if no permission for existing directory
func createOutputDirectory(dir string) error {
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			// if the dir already exists then test the write permission
			testFile := filepath.Join(dir, ".perm_check")
			f, writeErr := os.Create(testFile)
			if writeErr != nil {
				return fmt.Errorf("no write permission in %q: %w", dir, writeErr)
			}
			err := f.Close()
			if err != nil {
				return err
			}
			err = os.Remove(testFile) // clean up
			return err
		}
		// i think parent dir permission also go here? not sure
		return err
	}

	return nil
}

func generateStandardTemplate(fileTemplate *assets.FileTemplate, templateName, outDir string, data any) error {
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

	defer func() {
		if closeError := file.Close(); closeError != nil {
			if err == nil {
				err = closeError
			}
		}
	}()

	// execute template with the data needed
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("error executing template %s: %w", templateName, err)
	}

	return nil
}

func copyFile(src, dst string) error {
	err := assets.CopyEmbeddedFiles(src, dst)
	return err
}

func copyDir(src, dst string) error {
	err := assets.CopyEmbeddedDirectory(src, dst)
	return err
}
