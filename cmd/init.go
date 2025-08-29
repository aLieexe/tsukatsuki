/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/aLieexe/tsukatsuki/internal/ui/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type UserInput struct {
	ProductName *textinput.Output
	// ServerIP
	// ProductDomain
	// Webserver
	// Services
	// Languages
	// GithubActions
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a configuration file",
	Long:  `This should create a configuration file, that later can be used to deploy using tsukatsuki deploy`,
	Run: func(cmd *cobra.Command, args []string) {
		userInput := &UserInput{
			ProductName: &textinput.Output{},
		}

		initModel := textinput.InitializeTextinputModel(userInput.ProductName)

		program := tea.NewProgram(initModel)
		if _, err := program.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(userInput.ProductName.Value)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
