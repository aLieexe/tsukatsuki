package services

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/aLieexe/tsukatsuki/internal/assets"
	"github.com/aLieexe/tsukatsuki/internal/utils"
)

// List out all compose services to add in docker-compose.yaml
type ComposeConfig struct {
	Storage  []string
	Services []string
}

func (app *AppConfig) GenerateDeploymentFiles() error {
	// Pre configuration, adjust needed generation config etc
	extraRoles := make([]string, 0)

	extraRoles = append(extraRoles, "common", "docker")
	for _, actions := range app.GithubActions {
		if actions.Type == "actions-cd" {
			app.IPv6 = utils.IsIPv6(app.ServerIP)
			extraRoles = append(extraRoles, "cd-setup")
		}
	}

	extraConfigFiles := []string{fmt.Sprintf("dockerfile-%s", app.Runtime), "rsync-ignore"}

	operations := []struct {
		name string
		fn   func() error
	}{
		// This should be fine for now, can add the service later
		{"project compose generation", func() error {
			return app.GenerateProjectCompose()
		}},
		{"ansible files generation", func() error {
			return app.GenerateAnsibleFiles(extraRoles)
		}},
		{"configuration files generation", func() error {
			return app.GenerateConfigurationFiles(extraConfigFiles)
		}},
		{"proxy files generation", func() error {
			return app.GenerateProxyFiles()
		}},
		{"github actions files generation", func() error {
			return app.GenerateActionsFiles()
		}},
	}

	for _, op := range operations {
		if err := op.fn(); err != nil {
			return fmt.Errorf("%s : %w", op.name, err)
		}
	}

	return nil
}

func (app *AppConfig) GenerateProxyFiles() error {
	const composeTemplateName = "compose-proxy"
	proxyPresetCode := fmt.Sprintf("%s-proxy", app.Proxy)
	templateProvider, err := assets.NewTemplateProvider(app.OutputDir, app.ProjectName)
	if err != nil {
		return err
	}
	presetProvider := templateProvider.GetComposePresetTemplates()
	preset := presetProvider[proxyPresetCode]

	proxyPresetTmpl, err := template.New(proxyPresetCode).Option("missingkey=error").Parse(string(preset.Content))
	if err != nil {
		return fmt.Errorf("parsing template %s: %w", composeTemplateName, err)
	}

	presetTemplateData := struct {
		DockerImage string
		ACMEEmail   string
	}{
		DockerImage: app.ProxyImage,
		ACMEEmail:   app.ACMEEmail,
	}

	var buffer bytes.Buffer
	err = proxyPresetTmpl.Execute(&buffer, presetTemplateData)
	if err != nil {
		return fmt.Errorf("executing template %s: %w", proxyPresetCode, err)
	}

	// All the service and volumes listed previously
	proxyPreset := string(buffer.String())

	composeTemplateData := struct {
		Volumes       []string
		ProxyTemplate string
		Proxy         string
	}{
		Volumes:       preset.Volume,
		ProxyTemplate: proxyPreset,
		Proxy:         app.Proxy,
	}

	proxyComposeTemplate := templateProvider.GetFileTemplates()[composeTemplateName]

	err = createOutputDirectory(proxyComposeTemplate.OutputDir)
	if err != nil {
		return err
	}

	tmpl, err := template.New(composeTemplateName).Option("missingkey=error").Parse(string(proxyComposeTemplate.Content))
	if err != nil {
		return fmt.Errorf("parsing template %s: %w", proxyComposeTemplate.Filename, err)
	}

	// create output file
	filePath := filepath.Join(proxyComposeTemplate.OutputDir, proxyComposeTemplate.Filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", filePath, err)
	}

	defer func() {
		if closeError := file.Close(); closeError != nil {
			if err == nil {
				err = closeError
			}
		}
	}()

	err = tmpl.Execute(file, composeTemplateData)
	if err != nil {
		return fmt.Errorf("executing template %s: %w", composeTemplateName, err)
	}

	// Proxy don't need other files, only a docker compose file
	if app.Proxy == "traefik" {
		return nil
	}

	proxyEntryCode := fmt.Sprintf("proxy-%s-entry", app.Proxy)
	proxyProjectCode := fmt.Sprintf("proxy-%s-project", app.Proxy)

	fileTemplate := templateProvider.GetFileTemplates()[proxyEntryCode]
	if err := generateStandardTemplate(&fileTemplate, proxyEntryCode, app); err != nil {
		return fmt.Errorf("generating template %s: %w", proxyEntryCode, err)
	}

	fileTemplate = templateProvider.GetFileTemplates()[proxyProjectCode]

	err = createOutputDirectory(fileTemplate.OutputDir)
	if err != nil {
		return err
	}

	if err := generateStandardTemplate(&fileTemplate, proxyProjectCode, app); err != nil {
		return fmt.Errorf("generating template %s: %w", proxyProjectCode, err)
	}
	return nil
}

// TODO: Refactor it to be able to do array, Planning on changing the github stuff into multi choice instead
func (app *AppConfig) GenerateActionsFiles() error {
	if app.GithubActions == nil {
		return nil
	}

	templateProvider, err := assets.NewTemplateProvider(app.OutputDir, app.ProjectName)
	if err != nil {
		return err
	}

	for _, actions := range app.GithubActions {
		if actions.Type == "actions-cd" {
			code := fmt.Sprintf("%s-%s", actions.Type, "docker")
			fileTemplate := templateProvider.GetFileTemplates()[code]

			if err := generateStandardTemplate(&fileTemplate, code, app); err != nil {
				return fmt.Errorf("generating template %s: %w", code, err)
			}
			continue
		}
		code := fmt.Sprintf("%s-%s", actions.Type, app.Runtime)
		fileTemplate := templateProvider.GetFileTemplates()[code]

		if err := generateStandardTemplate(&fileTemplate, code, app); err != nil {
			return fmt.Errorf("generating template %s: %w", code, err)
		}
	}

	return nil
}

func (app *AppConfig) GenerateAnsibleFiles(extraRoles []string) error {
	const ansibleSetupCode = "ansible-setup"
	const ansibleInventoryCode = "ansible-inventory"
	const ansibleVarsCode = "ansible-vars"
	const ansibleProxyRolesCode = "ansible-proxy"

	ansibleStaticDir := "static/ansible"

	playbookData := struct {
		Roles []string
	}{
		Roles: make([]string, 0),
	}

	playbookData.Roles = append(playbookData.Roles, extraRoles...)

	if app.Security {
		playbookData.Roles = append(playbookData.Roles, "security")
	}

	templateProvider, err := assets.NewTemplateProvider(app.OutputDir, app.ProjectName)
	if err != nil {
		return fmt.Errorf("initializing template provider %s: %w", app.OutputDir, err)
	}

	fileTemplate := templateProvider.GetFileTemplates()[ansibleSetupCode]
	if err := generateStandardTemplate(&fileTemplate, ansibleSetupCode, playbookData); err != nil {
		return fmt.Errorf("generating template %s: %w", ansibleSetupCode, err)
	}

	fileTemplate = templateProvider.GetFileTemplates()[ansibleInventoryCode]
	if err := generateStandardTemplate(&fileTemplate, ansibleInventoryCode, app); err != nil {
		return fmt.Errorf("generating template %s: %w", ansibleInventoryCode, err)
	}

	fileTemplate = templateProvider.GetFileTemplates()[ansibleVarsCode]
	if err := generateStandardTemplate(&fileTemplate, ansibleVarsCode, app); err != nil {
		return fmt.Errorf("generating template %s: %w", ansibleVarsCode, err)
	}

	fileTemplate = templateProvider.GetFileTemplates()[ansibleProxyRolesCode]
	if err := generateStandardTemplate(&fileTemplate, ansibleProxyRolesCode, app); err != nil {
		return fmt.Errorf("generating template %s: %w", ansibleProxyRolesCode, err)
	}

	staticProvider := assets.NewStaticProvider(app.OutputDir)
	staticFile := staticProvider.StaticFile["ansible-config"]
	err = copyFile(staticFile.StaticFilePath, staticFile.OutputPath)
	if err != nil {
		return err
	}

	staticFile = staticProvider.StaticFile["ansible-deploy"]
	err = copyFile(staticFile.StaticFilePath, staticFile.OutputPath)
	if err != nil {
		return err
	}

	// Roles
	rolesSrcDir := filepath.Join(ansibleStaticDir, "/roles")
	rolesDstDir := filepath.Join(app.OutputDir, "ansible", "/roles")

	playbookData.Roles = append(playbookData.Roles, "deployment")

	for _, role := range playbookData.Roles {
		src := filepath.Join(rolesSrcDir, role)
		dst := filepath.Join(rolesDstDir, role)

		if err := copyDir(src, dst); err != nil {
			return fmt.Errorf("copying %s, in %s to %s: %w", role, src, dst, err)
		}
	}

	return nil
}

func (app *AppConfig) GenerateConfigurationFiles(templateNeeded []string) error {
	templateProvider, err := assets.NewTemplateProvider(app.OutputDir, app.ProjectName)
	if err != nil {
		return err
	}

	for _, templateName := range templateNeeded {
		fileTemplate := templateProvider.GetFileTemplates()[templateName]
		if err := generateStandardTemplate(&fileTemplate, templateName, app); err != nil {
			return fmt.Errorf("generating template %s: %w", templateName, err)
		}
	}
	return nil
}

func (app *AppConfig) GenerateProjectCompose() error {
	// Mapping name of project-compose.yaml in template_provider.go
	const composeTemplateName = "compose-project"

	templateProvider, err := assets.NewTemplateProvider(app.OutputDir, app.ProjectName)
	if err != nil {
		return err
	}
	composeTemplate := templateProvider.GetFileTemplates()[composeTemplateName]

	err = createOutputDirectory(composeTemplate.OutputDir)
	if err != nil {
		return err
	}

	tmpl, err := template.New(composeTemplateName).Option("missingkey=error").Parse(string(composeTemplate.Content))
	if err != nil {
		return fmt.Errorf("parsing template %s: %w", composeTemplate.Filename, err)
	}

	// create output file
	filePath := filepath.Join(composeTemplate.OutputDir, composeTemplate.Filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", filePath, err)
	}

	defer func() {
		if closeError := file.Close(); closeError != nil {
			if err == nil {
				err = closeError
			}
		}
	}()

	// Combine all the needed data, that is the services and the volumes needed for said service to function
	templateData := struct {
		Service        []string
		Volumes        []string
		ProjectName    string
		AppPort        int
		Proxy          string
		AppSiteAddress string
		Security       bool
	}{
		Service:        []string{},
		Volumes:        []string{},
		ProjectName:    app.ProjectName,
		AppPort:        app.AppPort,
		Proxy:          app.Proxy,
		AppSiteAddress: app.AppSiteAddress,
		Security:       app.Security,
	}

	// Combine services and proxy, why do i seperate this again?
	services := []Service{}

	for _, service := range app.Services {
		services = append(services, Service{
			Name:        service.Name,
			DockerImage: service.DockerImage,
		})
	}

	presetProvider := templateProvider.GetComposePresetTemplates()
	for _, service := range services {
		if preset, exists := presetProvider[service.Name]; exists {
			// Exec all the preset byitself
			serviceTmpl, err := template.New(service.Name).Option("missingkey=error").Parse(string(preset.Content))
			if err != nil {
				return fmt.Errorf("parsing template %s: %w", service.Name, err)
			}

			var buffer bytes.Buffer
			err = serviceTmpl.Execute(&buffer, service)
			if err != nil {
				return fmt.Errorf("executing template %s: %w", service.Name, err)
			}

			// All the service and volumes listed previously
			serviceDefinition := string(buffer.String())

			templateData.Service = append(templateData.Service, serviceDefinition)

			if preset.Volume != nil {
				templateData.Volumes = append(templateData.Volumes, preset.Volume...)
			}
		}
	}

	err = tmpl.Execute(file, templateData)
	if err != nil {
		return fmt.Errorf("executing template %s: %w", composeTemplateName, err)
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

			if closeErr := f.Close(); closeErr != nil {
				return fmt.Errorf("closing test file %s: %w", f.Name(), closeErr)
			}

			// Clean up test
			if removeErr := os.Remove(testFile); removeErr != nil {
				return fmt.Errorf("removing test file %s: %w", f.Name(), removeErr)
			}
			return nil
		}
		// i think parent dir permission also go here? not sure
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	return nil
}

func generateStandardTemplate(fileTemplate *assets.FileTemplate, templateName string, data any) error {
	err := createOutputDirectory(fileTemplate.OutputDir)
	if err != nil {
		return err
	}

	content := string(fileTemplate.Content)
	if content == "" {
		return fmt.Errorf("template content is empty for %s", templateName)
	}

	tmpl, err := template.New(templateName).Option("missingkey=error").Parse(content)
	if err != nil {
		return fmt.Errorf("parsing template '%s': %w", templateName, err)
	}

	// create output file
	filePath := filepath.Join(fileTemplate.OutputDir, fileTemplate.Filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", filePath, err)
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
		return fmt.Errorf("executing template %s: %w", templateName, err)
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
