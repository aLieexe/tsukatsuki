/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"

	"github.com/aLieexe/tsukatsuki/internal/config"
	"github.com/aLieexe/tsukatsuki/internal/log"
	"github.com/aLieexe/tsukatsuki/internal/services"
	"github.com/aLieexe/tsukatsuki/internal/ui/textinput"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Connect to remote machine via SSH",
	Long: `Connect to the remote machine with configuration defined in tsukatsuki.yaml.
You will be connected as the setup user`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.InitLogger(cmd)

		if !config.ConfigFileExist() {
			logger.Warn("Please generate a tsukatsuki.yaml with tsukatsuki init before continuing")
			os.Exit(1)
		}

		yamlConfig, err := config.GetConfigFromFiles()
		if err != nil {
			if err != nil {
				logger.Error(fmt.Sprintf("failed to read config from files: %s", err))
				os.Exit(1)
			}
		}

		app := services.NewAppConfigFromYaml(yamlConfig)
		var sshConfig *ssh.ClientConfig

		port, err := services.ProbeSSH(app.ServerIP, []int{app.SSHPort, 22})
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to find open SSH Port: %s", err))
		}
		logger.Info(fmt.Sprintf("Using SSH port %d for connection", port))

		keyPath := filepath.Join(app.OutputDir, "ansible", "key", app.SetupUser)
		key, err := os.ReadFile(keyPath)
		if err != nil {
			logger.Warn(fmt.Sprintf(`Failed reading key file in "%s" : %s`, keyPath, err))
			logger.Warn("Falling back to use password instead")

			sshConfig = getSSHConfigWithPassword(app, logger)
			err = services.ConnectSSH(sshConfig, app.ServerIP, port)
			if err != nil {
				logger.Error(fmt.Sprintf("Connection failed: %s", err))
				os.Exit(1)
			}
			os.Exit(1)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			var passErr *ssh.PassphraseMissingError
			if errors.As(err, &passErr) {
				sshConfig = getSSHConfigWithPassphrase(app, logger, keyPath)
			} else {
				logger.Error(fmt.Sprintf("Failed parsing private key: %s", err))
				os.Exit(1)
			}
		} else {
			sshConfig = &ssh.ClientConfig{
				User:            app.SetupUser,
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				Auth: []ssh.AuthMethod{
					ssh.PublicKeys(signer),
				},
			}
		}

		err = services.ConnectSSH(sshConfig, app.ServerIP, port)
		if err != nil {
			logger.Error(fmt.Sprintf("Connection failed: %s", err))
			os.Exit(1)
		}
	},
}

func getSSHConfigWithPassphrase(app *services.AppConfig, logger *slog.Logger, keyPath string) *ssh.ClientConfig {
	var passphrase textinput.Output

	teaProgram := tea.NewProgram(textinput.InitializePasswordInputModel(&passphrase, fmt.Sprintf(`What is the passphrase for the key located in %s`, keyPath), "12345678", app))
	_, err := teaProgram.Run()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to receive input %s: ", err))
		os.Exit(1)
	}

	key, err := os.ReadFile(keyPath)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed reading file: %s", err))
		os.Exit(1)
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(passphrase.Value))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed parsing private key: %s", err))
		os.Exit(1)
	}
	sshConfig := &ssh.ClientConfig{
		User:            app.SetupUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	return sshConfig
}

func getSSHConfigWithPassword(app *services.AppConfig, logger *slog.Logger) *ssh.ClientConfig {
	var password textinput.Output

	teaProgram := tea.NewProgram(textinput.InitializePasswordInputModel(&password, fmt.Sprintf(`What is the password for user "%s"`, app.SetupUser), "12345678", app))
	_, err := teaProgram.Run()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to receive input %s: ", err))
		os.Exit(1)
	}

	sshConfig := &ssh.ClientConfig{
		User:            app.SetupUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(password.Value),
		},
	}

	return sshConfig
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
