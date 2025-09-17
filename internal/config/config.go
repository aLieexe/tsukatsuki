package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type AppConfigYaml struct {
	Project struct {
		Name    string `yaml:"name"`
		Domain  string `yaml:"domain"`
		Port    int    `yaml:"port"`
		Runtime string `yaml:"runtime"`
	} `yaml:"project"`

	Server struct {
		IP string `yaml:"ip"`
	} `yaml:"server"`

	Webserver struct {
		Type     string `yaml:"type"`
		SSLEmail string `yaml:"ssl_email"`
	} `yaml:"webserver"`

	// Services []struct {
	// }

	GithubActions struct {
		Mode   string `yaml:"mode"`
		Branch string `yaml:"branch"`
	} `yaml:"github_actions"`
}

func CreateConfigFiles(cfg AppConfigYaml) error {
	file, err := os.ReadFile("tsukatsuki.yaml")
	if err != nil {
		log.Println("tsukatsuki.yaml file doesnt exist in the current project, re-creating one", err)
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

	err = os.WriteFile("tsukatsuki.yaml", modifiedData, 0644)
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
