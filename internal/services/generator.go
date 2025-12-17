package services

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/aLieexe/tsukatsuki/internal/assets"
	"github.com/aLieexe/tsukatsuki/internal/utils"
)

type ComposeConfig struct {
	Storage  []string
	Services []string
}

func (app *AppConfig) GenerateDeploymentFiles(logger *slog.Logger) error {
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
		fn   func(*slog.Logger) error
	}{
		{"project compose generation", func(l *slog.Logger) error {
			return app.GenerateProjectCompose(l)
		}},
		{"compose helper generation", func(l *slog.Logger) error {
			return app.GenerateHelperCompose(l)
		}},
		{"ansible files generation", func(l *slog.Logger) error {
			return app.GenerateAnsibleFiles(extraRoles, l)
		}},
		{"configuration files generation", func(l *slog.Logger) error {
			return app.GenerateConfigurationFiles(extraConfigFiles, l)
		}},
		{"proxy files generation", func(l *slog.Logger) error {
			return app.GenerateProxyFiles(l)
		}},
		{"github actions files generation", func(l *slog.Logger) error {
			return app.GenerateActionsFiles(l)
		}},
	}

	for _, op := range operations {
		if err := op.fn(logger); err != nil {
			return fmt.Errorf("%s : %w", op.name, err)
		}
	}

	return nil
}

// TODO
func (app *AppConfig) GenerateHelperCompose(logger *slog.Logger) error {
	return nil
}

func (app *AppConfig) GenerateProxyFiles(logger *slog.Logger) error {
	// Proxy compose file code based on the one defined on template_provider.go
	const composeTemplateName = "compose-proxy"

	templateProvider, err := assets.NewTemplateProvider(app.OutputDir, app.ProjectName)
	if err != nil {
		return err
	}

	proxyComposeTemplate := templateProvider.GetFileTemplates()[composeTemplateName]

	filePath := filepath.Join(proxyComposeTemplate.OutputDir, proxyComposeTemplate.Filename)

	// If a proxy compose file already existed, skip it
	if _, err := os.Stat(filePath); err == nil {
		logger.Warn(fmt.Sprintf("skipping %s: file already exists", filePath))
		if app.Proxy != "traefik" {
			proxyEntryCode := fmt.Sprintf("proxy-%s-entry", app.Proxy)
			proxyProjectCode := fmt.Sprintf("proxy-%s-project", app.Proxy)

			fileTemplate := templateProvider.GetFileTemplates()[proxyEntryCode]
			if err := generateStandardTemplate(&fileTemplate, proxyEntryCode, app, logger); err != nil {
				return fmt.Errorf("generating template %s: %w", proxyEntryCode, err)
			}

			fileTemplate = templateProvider.GetFileTemplates()[proxyProjectCode]
			err = createOutputDirectory(fileTemplate.OutputDir)
			if err != nil {
				return err
			}

			if err := generateStandardTemplate(&fileTemplate, proxyProjectCode, app, logger); err != nil {
				return fmt.Errorf("generating template %s: %w", proxyProjectCode, err)
			}
		}
		return nil
	}

	// Get the specific template for each proxy that will be used within the proxy compose file
	proxyPresetCode := fmt.Sprintf("%s-proxy", app.Proxy)
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

	// The proxy preset template generated above
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

	err = createOutputDirectory(proxyComposeTemplate.OutputDir)
	if err != nil {
		return err
	}

	tmpl, err := template.New(composeTemplateName).Option("missingkey=error").Parse(string(proxyComposeTemplate.Content))
	if err != nil {
		return fmt.Errorf("parsing template %s: %w", proxyComposeTemplate.Filename, err)
	}

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

	// Some proxy don't need other files, only a docker-compose file
	// Traefik is set to use its docker provider, not needing other file to functions
	if app.Proxy == "traefik" {
		return nil
	}

	// A code referring to the proxy global config that will import the project specific file
	// Both code are defined within template_provider.go
	proxyEntryCode := fmt.Sprintf("proxy-%s-entry", app.Proxy)
	// Code referring to the project specific proxy file configuration
	proxyProjectCode := fmt.Sprintf("proxy-%s-project", app.Proxy)

	fileTemplate := templateProvider.GetFileTemplates()[proxyEntryCode]
	if err := generateStandardTemplate(&fileTemplate, proxyEntryCode, app, logger); err != nil {
		return fmt.Errorf("generating template %s: %w", proxyEntryCode, err)
	}

	fileTemplate = templateProvider.GetFileTemplates()[proxyProjectCode]

	err = createOutputDirectory(fileTemplate.OutputDir)
	if err != nil {
		return err
	}

	if err := generateStandardTemplate(&fileTemplate, proxyProjectCode, app, logger); err != nil {
		return fmt.Errorf("generating template %s: %w", proxyProjectCode, err)
	}
	return nil
}

func (app *AppConfig) GenerateActionsFiles(logger *slog.Logger) error {
	if app.GithubActions == nil {
		return nil
	}

	templateProvider, err := assets.NewTemplateProvider(app.OutputDir, app.ProjectName)
	if err != nil {
		return err
	}

	for _, actions := range app.GithubActions {
		// Currently done this way for future expansion
		// Since docker is used there wont be any difference between command
		// Which is why they end up using the same command
		if actions.Type == "actions-cd" {
			code := fmt.Sprintf("%s-%s", actions.Type, "docker")
			fileTemplate := templateProvider.GetFileTemplates()[code]

			if err := generateStandardTemplate(&fileTemplate, code, app, logger); err != nil {
				return fmt.Errorf("generating template %s: %w", code, err)
			}
			continue
		}
		code := fmt.Sprintf("%s-%s", actions.Type, app.Runtime)
		fileTemplate := templateProvider.GetFileTemplates()[code]

		// Use a generic function to generate the file
		if err := generateStandardTemplate(&fileTemplate, code, app, logger); err != nil {
			return fmt.Errorf("generating template %s: %w", code, err)
		}
	}

	return nil
}

func (app *AppConfig) GenerateAnsibleFiles(extraRoles []string, logger *slog.Logger) error {
	// Ansible files code that are specified in template_provider.go
	const ansibleSetupCode = "ansible-setup"
	const ansibleInventoryCode = "ansible-inventory"
	const ansibleVarsCode = "ansible-vars"
	const ansibleProxyRolesCode = "ansible-proxy"

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
	if err := generateStandardTemplate(&fileTemplate, ansibleSetupCode, playbookData, logger); err != nil {
		return fmt.Errorf("generating template %s: %w", ansibleSetupCode, err)
	}

	fileTemplate = templateProvider.GetFileTemplates()[ansibleInventoryCode]
	if err := generateStandardTemplate(&fileTemplate, ansibleInventoryCode, app, logger); err != nil {
		return fmt.Errorf("generating template %s: %w", ansibleInventoryCode, err)
	}

	fileTemplate = templateProvider.GetFileTemplates()[ansibleVarsCode]
	if err := generateStandardTemplate(&fileTemplate, ansibleVarsCode, app, logger); err != nil {
		return fmt.Errorf("generating template %s: %w", ansibleVarsCode, err)
	}

	fileTemplate = templateProvider.GetFileTemplates()[ansibleProxyRolesCode]
	if err := generateStandardTemplate(&fileTemplate, ansibleProxyRolesCode, app, logger); err != nil {
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
	ansibleStaticDir := "static/ansible"
	rolesSrcDir := filepath.Join(ansibleStaticDir, "/roles")
	rolesDstDir := filepath.Join(app.OutputDir, "ansible", "/roles")

	playbookData.Roles = append(playbookData.Roles, "deployment")

	// Copy each roles one by one
	for _, role := range playbookData.Roles {
		src := filepath.Join(rolesSrcDir, role)
		dst := filepath.Join(rolesDstDir, role)

		if err := copyDir(src, dst); err != nil {
			return fmt.Errorf("copying %s, in %s to %s: %w", role, src, dst, err)
		}
	}

	return nil
}

func (app *AppConfig) GenerateConfigurationFiles(templateNeeded []string, logger *slog.Logger) error {
	templateProvider, err := assets.NewTemplateProvider(app.OutputDir, app.ProjectName)
	if err != nil {
		return err
	}

	for _, templateName := range templateNeeded {
		// Use a generic function to generate the file
		fileTemplate := templateProvider.GetFileTemplates()[templateName]
		if err := generateStandardTemplate(&fileTemplate, templateName, app, logger); err != nil {
			return fmt.Errorf("generating template %s: %w", templateName, err)
		}
	}
	err = ensureGitignoreExist()
	if err != nil {
		return fmt.Errorf("ensuring .gitignore: %w", err)
	}
	return nil
}

// An idempotent function that ensure some file is not commited
func ensureGitignoreExist() error {
	existingContent := make(map[string]int)
	filename := ".gitignore"
	contents := []string{"key", ".ansible", "app.tar"}

	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		existingContent[line] = 0
	}

	defer func() {
		if closeError := file.Close(); closeError != nil {
			if err == nil {
				err = closeError
			}
		}
	}()

	file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	defer func() {
		if closeError := file.Close(); closeError != nil {
			if err == nil {
				err = closeError
			}
		}
	}()

	for _, content := range contents {
		if _, found := existingContent[content]; !found {
			if _, err := file.WriteString(content + "\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

func (app *AppConfig) GenerateProjectCompose(logger *slog.Logger) error {
	// Mapping name of project-compose.yaml in template_provider.go
	const composeTemplateName = "compose-project"

	templateProvider, err := assets.NewTemplateProvider(app.OutputDir, app.ProjectName)
	if err != nil {
		return err
	}
	projectComposeTemplate := templateProvider.GetFileTemplates()[composeTemplateName]

	// If a project compose file already existed, skip it
	filePath := filepath.Join(projectComposeTemplate.OutputDir, projectComposeTemplate.Filename)
	if _, err := os.Stat(filePath); err == nil {
		logger.Warn(fmt.Sprintf("skipping %s: file already exists", filePath))
		return nil
	}

	err = createOutputDirectory(projectComposeTemplate.OutputDir)
	if err != nil {
		return err
	}

	tmpl, err := template.New(composeTemplateName).Option("missingkey=error").Parse(string(projectComposeTemplate.Content))
	if err != nil {
		return fmt.Errorf("parsing template %s: %w", projectComposeTemplate.Filename, err)
	}

	// create output file
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
	// This will be used in executing the project-compose file template
	templateData := struct {
		ServiceName     []string
		ServiceTemplate []string
		Volumes         []string
		ProjectName     string
		AppPort         int
		Proxy           string
		AppSiteAddress  string
		Security        bool
		EnvFile         string
	}{
		ServiceTemplate: []string{},
		Volumes:         []string{},
		ProjectName:     app.ProjectName,
		AppPort:         app.AppPort,
		Proxy:           app.Proxy,
		AppSiteAddress:  app.AppSiteAddress,
		Security:        app.Security,
		EnvFile:         app.EnvFile,
	}

	for _, service := range app.Services {
		templateData.ServiceName = append(templateData.ServiceName, service.Name)
	}

	presetProvider := templateProvider.GetComposePresetTemplates()

	// The env is here rather then on its own loop for better performance (honestly, it's either negligible or worse, idk why i do this)
	env := []assets.EnvVar{}

	// Prepare each services preset that will be included in projec compose template
	for _, service := range app.Services {
		if preset, exists := presetProvider[service.Name]; exists {
			// Exec all the preset byitself
			serviceTmpl, err := template.New(service.Name).Option("missingkey=error").Parse(string(preset.Content))
			if err != nil {
				return fmt.Errorf("parsing template %s: %w", service.Name, err)
			}

			data := struct {
				DockerImage string
				Name        string
				ProjectName string
			}{
				Name:        service.Name,
				DockerImage: service.DockerImage,
				ProjectName: app.ProjectName,
			}

			var buffer bytes.Buffer
			err = serviceTmpl.Execute(&buffer, data)
			if err != nil {
				return fmt.Errorf("executing template %s: %w", service.Name, err)
			}

			// All the service and volumes listed previously
			serviceDefinition := string(buffer.String())

			templateData.ServiceTemplate = append(templateData.ServiceTemplate, serviceDefinition)

			if preset.Volume != nil {
				for i, volume := range preset.Volume {
					preset.Volume[i] = fmt.Sprintf("%s_%s", app.ProjectName, volume)
				}
				templateData.Volumes = append(templateData.Volumes, preset.Volume...)
			}

			env = append(env, preset.EnvVar...)

		}
	}

	err = tmpl.Execute(file, templateData)
	if err != nil {
		return fmt.Errorf("executing template %s: %w", composeTemplateName, err)
	}

	err = ensureEnvVars(app.EnvFile, env)
	if err != nil {
		return fmt.Errorf("ensuring env vars exist: %w", err)
	}
	return nil
}

// Ensure that .env file have a placeholder, these placeholder is needed for the services to function
// This will also clear up and give example for the user. They are defined in template_provider.go
func ensureEnvVars(path string, vars []assets.EnvVar) error {
	existingVars := make(map[string]int)

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if k, _, ok := strings.Cut(line, "="); ok {
			existingVars[k] = 0
		}
	}

	defer func() {
		if closeError := file.Close(); closeError != nil {
			if err == nil {
				err = closeError
			}
		}
	}()

	file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	defer func() {
		if closeError := file.Close(); closeError != nil {
			if err == nil {
				err = closeError
			}
		}
	}()

	for _, v := range vars {
		if _, found := existingVars[v.Name]; !found {
			if _, err := file.WriteString(v.Name + "=" + v.Default + "\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

// Create output directory, if not exist
// return error if no permission for existing directory
func createOutputDirectory(dir string) error {
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}
	return nil
}

func generateStandardTemplate(fileTemplate *assets.FileTemplate, templateName string, data any, logger *slog.Logger) error {
	err := createOutputDirectory(fileTemplate.OutputDir)
	if err != nil {
		return err
	}

	// If file already existed, skip the current file
	filePath := filepath.Join(fileTemplate.OutputDir, fileTemplate.Filename)
	if _, err := os.Stat(filePath); err == nil {
		logger.Warn(fmt.Sprintf("skipping %s: file already exists", filePath))
		return nil
	}

	content := string(fileTemplate.Content)
	if content == "" {
		return fmt.Errorf("template content is empty for %s", templateName)
	}

	tmpl, err := template.New(templateName).Option("missingkey=error").Parse(content)
	if err != nil {
		return fmt.Errorf("parsing template '%s': %w", templateName, err)
	}

	// Create output file
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

	// Execute template with the data needed
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
