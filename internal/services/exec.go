package services

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/x/term"
	"golang.org/x/crypto/ssh"

	"github.com/aLieexe/tsukatsuki/internal/utils"
)

func checkTCPReachable(host string, port int, timeout time.Duration) error {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return fmt.Errorf("dialing TCP %s: %w", addr, err)
	}

	// I don't think this ever fail, but uh yes
	err = conn.Close()
	if err != nil {
		return fmt.Errorf("closing TCP connection to %s: %w", addr, err)
	}

	return nil
}

func ProbeSSH(host string, portList []int) (int, error) {
	var err error
	for _, port := range portList {
		err = checkTCPReachable(host, port, 5*time.Second)
		if err == nil {
			return port, nil
		}
	}
	return 0, err
}

func execCmd(cmd *exec.Cmd, logger *slog.Logger, errorPatterns ...string) error {
	var stdoutBuf, stderrBuf bytes.Buffer

	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	stdout := stdoutBuf.String()
	if err != nil {
		for _, pattern := range errorPatterns {
			if strings.Contains(stdout, pattern) {
				return fmt.Errorf("playbook error: %s", pattern)
			}
		}

		if strings.Contains(stdout, "PLAY RECAP") {
			recapSection := strings.Split(stdout, "PLAY RECAP")[1]
			if strings.Contains(recapSection, "failed=") && !strings.Contains(recapSection, "failed=0") {
				return fmt.Errorf("playbook error: tasks failed")
			}
			if strings.Contains(recapSection, "unreachable=") && !strings.Contains(recapSection, "unreachable=0") {
				return fmt.Errorf("playbook error: hosts unreachable")
			}
		}

		return fmt.Errorf("process error : %w", err)
	}

	return nil
}

func ExecCommand(cmd *exec.Cmd) (string, error) {
	var stdoutBuf, stderrBuf bytes.Buffer

	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("executing command %s: %w: %s", cmd.String(), err, stderrBuf.String())
	}

	return stdoutBuf.String(), nil
}

// Should be use as a first attempt / option before trying out the one with password
func ExecAnsible(logger *slog.Logger, ansiblePath, playbookName string, port int) error {
	cmd := exec.Command(
		"ansible-playbook",
		playbookName,
		"-i", "inventory.ini",
		"-c", "ssh",
		"-e", fmt.Sprintf("ssh_port=%d", port),
	)

	cmd.Dir = ansiblePath

	err := execCmd(cmd, logger, "no hosts matched")
	if err != nil {
		return fmt.Errorf("executing with inventory file: %w", err)
	}
	return nil
}

// This should be used as a fallback, in the case that the one with inventory.ini dont work
func ExecAnsibleWithPassword(logger *slog.Logger, ansiblePath, playbookName, password string, port int) error {
	cmd := exec.Command(
		"ansible-playbook",
		playbookName,
		"-i", "inventory.ini",
		"-c", "ssh",
		// "-e", fmt.Sprintf("ansible_become_pass=%s ansible_password=%s", password, password),
		"-e", fmt.Sprintf("ansible_become_pass=%s ansible_password=%s ssh_port=%d", password, password, port),
	)

	cmd.Dir = ansiblePath
	err := execCmd(cmd, logger, "no hosts matched")
	if err != nil {
		return fmt.Errorf("executing with password: %w", err)
	}
	return nil
}

func ConnectSSH(config *ssh.ClientConfig, serverIp string, sshPort int) error {
	var sshHost string
	if utils.IsIPv6(serverIp) {
		sshHost = fmt.Sprintf("[%s]:%d", serverIp, sshPort)
	} else {
		sshHost = fmt.Sprintf("%s:%d", serverIp, sshPort)
	}

	client, err := ssh.Dial("tcp", sshHost, config)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if deferError := client.Close(); deferError != nil {
			if err == nil {
				err = deferError
			}
		}
	}()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if deferError := session.Close(); deferError != nil {
			if err == nil {
				err = deferError
			}
		}
	}()

	fd := os.Stdin.Fd() // fd is uintptr

	// Make terminal raw
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if deferError := term.Restore(fd, oldState); deferError != nil {
			if err == nil {
				err = deferError
			}
		}
	}()

	// Get terminal size (expects int)
	width, height, err := term.GetSize(fd)
	if err != nil {
		log.Fatal(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
		log.Fatal(err)
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		log.Fatal(err)
	}

	if err := session.Wait(); err != nil {
		return err
	}
	return nil
}
