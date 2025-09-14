package services

// List out all compose services to add in docker-compose.yaml
type ComposeConfig struct {
	Storage  []string
	Services []string
}

// func (app *AppConfig) GenerateAnsibleFiles() {
// 	// TODO: Generate Ansible-playbook, with roles already prepared and configured to the AppConfig

// }

func (app *AppConfig) GenerateDockerFiles() {
	// TODO: Generate docker compose files, and all of its needed configuration (Caddyfile, Dockerfile, etc)

	generator.
}
