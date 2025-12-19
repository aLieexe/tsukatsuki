## Getting Started
Tsukatsuki is an open-source command-line tool designed to simplify application deployment on self-hosted infrastructure. It generates Ansible roles from selected configurations and can also produce GitHub Actions workflows to support CI/CD pipelines. Docker is used as the deployment mechanism to ensure consistent and reproducible environments.

Tsukatsuki currently supports **Node.js** and **Go** runtimes. The project focuses on a **single-Dockerfile workflow**: one Dockerfile is used to generate a Docker Compose setup for deployment. If multiple services or Dockerfiles are required, each must be generated and initialized as a separate Tsukatsuki project. The Dockerfile is expected to be deployment-ready before initialization.

---
## Installation
Tsukatsuki is available to 

- Linux
- Windows via WSL

As long as your OS can run ansible as the control node, then tsukatsuki will work fine. 
You can however still benefits even without it, especially for generate the configuration files needed

Clone the repository and install the CLI:
```
git clone https://github.com/aLieexe/tsukatsuki.git
cd tsukatsuki
make install
```
After installation, the `tsukatsuki` command becomes available system-wide.

---
## Initializing a Project
To generate a deployment setup, run:
```
tsukatsuki init
```
An interactive prompt guides the configuration process. Based on the selected options, Tsukatsuki generates:

- A `tsukatsuki.yaml` configuration file  
- Ansible roles required for deployment  
- Supporting files needed for Docker-based deployment  

All selected configurations are stored in `tsukatsuki.yaml` and reused for future commands.

Only one Dockerfile is supported per project. Each Dockerfile results in a single Docker Compose configuration. Multiple Dockerfiles must be initialized in separate project directories.

---
## Deploying to a Server
After initialization is complete, deployment can be started with:

```
tsukatsuki deploy
```

This command uses the generated configuration and Ansible roles to deploy the application to the target server using Docker.

---
## Managing server
For server management, currently tsukatsuki present an easy way to access it via SSH by running:
```
tsukatsuki ssh
```

A way to manage server easier is planned and will be continued to be worked on in the future