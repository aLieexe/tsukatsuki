package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

	if logger.Enabled(context.Background(), slog.LevelDebug) {
		logger.Debug("playbook output", "stdout", stdout)
	}

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
