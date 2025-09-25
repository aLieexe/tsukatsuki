package services

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func execCmd(cmd *exec.Cmd, errorPatterns ...string) error {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	mStdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	mStderr := io.MultiWriter(os.Stderr, &stderrBuf)

	cmd.Stdout = mStdout
	cmd.Stderr = mStderr

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("process error : %w", err)

	}

	stdout := stdoutBuf.String()
	for _, pattern := range errorPatterns {
		if strings.Contains(stdout, pattern) {
			return fmt.Errorf("playbook error: %s", pattern)
		}
	}

	if strings.Contains(stdout, "PLAY RECAP") {
		recapSection := strings.Split(stdout, "PLAY RECAP")[1]
		if strings.Contains(recapSection, "failed=") && !strings.Contains(recapSection, "failed=0") {
			return fmt.Errorf("playbook error: tasks failed on some hosts")
		}
	}
	return nil
}

// Should be use as a first attempt / option before trying out the one with password
func ExecAnsible(playbookPath string) error {
	cmd := exec.Command(
		"ansible-playbook",
		playbookPath,
		"-i", "inventory.ini",
		"-c", "ssh",
	)

	err := execCmd(cmd, "no hosts matched")
	if err != nil {
		return err
	}
	return nil
}

// This should be used as a fallback, in the case that the one with inventory.ini dont work
func ExecAnsibleWithPassword(playbookPath string, password string) error {
	cmd := exec.Command(
		"ansible-playbook",
		playbookPath,
		"-i", "inventory.ini",
		"-c", "ssh",
		"-e", fmt.Sprintf("ansible_become_pass=%s ansible_password=%s", password, password),
	)

	err := execCmd(cmd, "no hosts matched")
	if err != nil {
		return err
	}
	return nil
}
