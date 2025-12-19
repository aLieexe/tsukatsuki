After initialization, tsukatsuki will generate a file with the name of `tsukatsuki.yaml`. This config file will be needed and used for a lot of tsukatsuki command, making it mandatory.

## YAML Structure
`tsukatsuki.yaml` example
```yaml
project:
  name: tsukatsuki-test
  port: 5500
  runtime: go
  build_image: golang:1.24.4-bookworm
  entrypoint: ./cmd/api/main.go
server:
  ip: 2407:6ac0:3:9d:abcd::1c7
  setup_user: setup
  ssh_port: 2222
  security: true
proxy:
  type: caddy
  domain: app.example.com
  docker_image: caddy:2.10.2-alpine
  acme_email: admin@example.com
services:
- name: postgresql
  docker_image: postgres:18.0-alpine
- name: redis
  docker_image: redis:8.2-alpine3.22
- name: rabbitmq
  docker_image: rabbitmq:4.2-management-alpine
- name: minio
  docker_image: minio/minio:RELEASE.2025-09-07T16-13-09Z
path:
  local_path: /home/alie/code/golang/tsukatsuki-test
  remote_path: /home/tsukatsuki/tsukatsuki
  output_dir: deploy
  env_file: .env
github_actions:
- type: actions-cd
- type: actions-ci
```

| Key                       | Description                                                                     | Available Values                           |
| ------------------------- | ------------------------------------------------------------------------------- | ------------------------------------------ |
| `project.name`            | Used for naming, identification, and grouping deployment resources.             | —                                          |
| `project.port`            | Port the application listens on inside the container.                           | —                                          |
| `project.runtime`         | Application runtime type that determines build and execution behavior.          | `go`, `node`                               |
| `project.build_image`     | Base Docker image used to build the application Docker image.                   | —                                          |
| `project.entrypoint`      | File or executable used to start the application.                               | —                                          |
| `server.ip`               | IP address of the target deployment server. Supports IPv4 and IPv6.             | —                                          |
| `server.setup_user`       | System user used for deployment and file ownership. Root user is not allowed.   | —                                          |
| `server.ssh_port`         | SSH port used for remote server access.                                         | —                                          |
| `server.security`         | When enabled, `tsukatsuki generate` produces additional security Ansible roles. | `true`, `false`                            |
| `proxy.type`              | Reverse proxy implementation used to expose the application.                    | `traefik`, `caddy`                         |
| `proxy.domain`            | Public domain assigned to the application.                                      | —                                          |
| `proxy.docker_image`      | Docker image used to run the reverse proxy.                                     | —                                          |
| `proxy.acme_email`        | Email address used for automatic TLS certificate registration.                  | —                                          |
| `services`                | List of dependent services required by the application.                         | —                                          |
| `services[].name`         | Logical service name used for identification and internal networking.           | `postgresql`, `redis`, `rabbitmq`, `minio` |
| `services[].docker_image` | Docker image used to run the service.                                           | —                                          |
| `path.local_path`         | Local filesystem path containing the project source.                            | —                                          |
| `path.remote_path`        | Remote filesystem path where the project is deployed.                           | —                                          |
| `path.output_dir`         | Directory where deployment artifacts are generated.                             | —                                          |
| `path.env_file`           | Environment variable file loaded at runtime.                                    | —                                          |
| `github_actions`          | List of GitHub Actions workflows enabled for the project.                       | —                                          |
| `github_actions[].type`   | Type of GitHub Actions workflow to be generated.                                | `actions-ci`, `actions-cd`                 |


## Generated .env
Some services will need environments to be able to work, this too can include the application that will be deployed. After generating configuration, be it via `init` or `generate` command.
This `.env` will be used in the `project-compose.yaml` that will be generated, and will need manual configuration from the user.