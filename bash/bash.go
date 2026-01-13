// Package bash provides utilities for executing shell commands.
package bash

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func GetDefaultShell() (string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "", fmt.Errorf("SHELL environment variable is not set")
	}
	return shell, nil
}

func RunCommand(command string) (string, error) {
	// Get shell from environment variable
	shell, err := GetDefaultShell()
	if err != nil {
		// In case SHELL is not set, default to zsh
		shell = "/bin/zsh"
	}

	cmd := exec.Command(shell, "-c", command)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return stdout.String(), nil
}
