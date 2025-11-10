package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type AppConfigYaml struct {
	Project struct {
		Name       string `yaml:"name"`
		Port       int    `yaml:"port"`
		Runtime    string `yaml:"runtime"`
		BuildImage string `yaml:"build_image"`
	} `yaml:"project"`

	Server struct {
		IP        string `yaml:"ip"`
		SetupUser string `yaml:"setup_user"`
		SSHPort   int    `yaml:"ssh_port"`
		Security  bool   `yaml:"security"`
	} `yaml:"server"`

	Webserver struct {
		Type        string `yaml:"type"`
		Domain      string `yaml:"domain"`
		DockerImage string `yaml:"docker_image"`
	} `yaml:"webserver"`

	Services []struct {
		Name        string `yaml:"name"`
		DockerImage string `yaml:"docker_image"`
	} `yaml:"services"`

	Path struct {
		LocalPath  string `yaml:"local_path"`
		RemotePath string `yaml:"remote_path"`
		OutputDir  string `yaml:"output_dir"`
	} `yaml:"path"`

	GithubActions []struct {
		Type string `yaml:"type"`
	} `yaml:"github_actions"`
}

func UpdateConfigFile(cfg AppConfigYaml) error {
	const configFileName = "tsukatsuki.yaml"

	file, err := os.ReadFile(configFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return writeConfigFile(configFileName, cfg, yaml.CommentMap{})
		}
		return fmt.Errorf("reading config file: %w", err)
	}

	// file exists, parse to extract comments
	commentMap := yaml.CommentMap{}
	err = yaml.UnmarshalWithOptions(file, &AppConfigYaml{}, yaml.Strict(), yaml.CommentToMap(commentMap))
	if err != nil {
		// if parsing fails, create new file without comments
		return writeConfigFile(configFileName, cfg, yaml.CommentMap{})
	}

	// write config with preserved comments
	return writeConfigFile(configFileName, cfg, commentMap)
}

func writeConfigFile(fileName string, cfg AppConfigYaml, commentMap yaml.CommentMap) error {
	modifiedData, err := yaml.MarshalWithOptions(
		cfg,
		yaml.WithComment(commentMap),
	)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	err = os.WriteFile(fileName, modifiedData, 0o644)
	if err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

func GetConfigFromFiles() (AppConfigYaml, error) {
	file, err := os.ReadFile("tsukatsuki.yaml")
	if err != nil {
		return AppConfigYaml{}, fmt.Errorf("reading config file: %w", err)
	}

	var yamlResult AppConfigYaml
	err = yaml.UnmarshalWithOptions(file, &yamlResult, yaml.Strict())
	if err != nil {
		return AppConfigYaml{}, fmt.Errorf("unmarshaling YAML: %w", err)
	}

	return yamlResult, nil
}

// supposed to be used only to check re init? why did i make this again?
func ConfigFileExist() bool {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := filepath.Join(dir, "tsukatsuki.yaml")
	_, err = os.Stat(path)
	return err == nil
}
