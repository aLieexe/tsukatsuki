package services

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

func execCmd(cmd *exec.Cmd, logger *slog.Logger, errorPatterns ...string) error {
	var stdoutBuf, stderrBuf bytes.Buffer

	mStdout := &stdoutBuf
	mStderr := &stderrBuf

	cmd.Stdout = mStdout
	cmd.Stderr = mStderr

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
func ExecAnsible(logger *slog.Logger, ansiblePath, playbookName string) error {
	cmd := exec.Command(
		"ansible-playbook",
		playbookName,
		"-i", "inventory.ini",
		"-c", "ssh",
	)

	cmd.Dir = ansiblePath

	err := execCmd(cmd, logger, "no hosts matched")
	if err != nil {
		return fmt.Errorf("executing with inventory file: %w", err)
	}
	return nil
}

// This should be used as a fallback, in the case that the one with inventory.ini dont work
func ExecAnsibleWithPassword(logger *slog.Logger, ansiblePath, playbookName, password string) error {
	cmd := exec.Command(
		"ansible-playbook",
		playbookName,
		"-i", "inventory.ini",
		"-c", "ssh",
		"-e", fmt.Sprintf("ansible_become_pass=%s ansible_password=%s", password, password),
	)

	cmd.Dir = ansiblePath
	err := execCmd(cmd, logger, "no hosts matched")
	if err != nil {
		return fmt.Errorf("executing with password: %w", err)
	}
	return nil
}
