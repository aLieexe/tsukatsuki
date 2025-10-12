package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type AppConfigYaml struct {
	Project struct {
		Name    string `yaml:"name"`
		Port    int    `yaml:"port"`
		Runtime string `yaml:"runtime"`
	} `yaml:"project"`

	Server struct {
		IP        string `yaml:"ip"`
		SetupUser string `yaml:"setup_user"`
	} `yaml:"server"`

	Webserver struct {
		Type   string `yaml:"type"`
		Domain string `yaml:"domain"`
	} `yaml:"webserver"`

	// Services []struct {
	// }

	Path struct {
		LocalPath  string `yaml:"local_path"`
		RemotePath string `yaml:"remote_path"`
		OutputDir  string `yaml:"output_dir"`
	} `yaml:"path"`

	GithubActions struct {
		Mode   string `yaml:"mode"`
		Branch string `yaml:"branch"`
	} `yaml:"github_actions"`
}

func CreateConfigFiles(cfg AppConfigYaml) error {
	file, err := os.ReadFile("tsukatsuki.yaml")
	if err != nil {
		fmt.Println("tsukatsuki.yaml file doesnt exist in the current project, re-creating one")
	}

	var yamlResult AppConfigYaml
	commentMap := yaml.CommentMap{}
	err = yaml.UnmarshalWithOptions(file, &yamlResult, yaml.Strict(), yaml.CommentToMap(commentMap))
	if err != nil {
		return fmt.Errorf("unmarshaling YAML: %w", err)
	}

	// write data with the comment aswell
	modifiedData, _ := yaml.MarshalWithOptions(
		cfg,
		yaml.WithComment(commentMap),
	)

	err = os.WriteFile("tsukatsuki.yaml", modifiedData, 0o644)
	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

func GetConfigFromFiles() (AppConfigYaml, error) {
	file, _ := os.ReadFile("tsukatsuki.yaml")

	var yamlResult AppConfigYaml
	err := yaml.UnmarshalWithOptions(file, &yamlResult, yaml.Strict())
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
