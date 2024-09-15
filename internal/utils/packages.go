package utils

import (
	"errors"
	"os/exec"
	"strings"
)

func CheckExists(binary string, versionCmd []string) (bool, string, error) {
	_, err := exec.LookPath(binary)
	if err != nil {
		return false, "", nil
	}

	cmd := exec.Command(versionCmd[0], versionCmd[1:]...)
	output, err := cmd.Output()

	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return true, "", errors.New(string(exitError.Stderr) + "\n" + exitError.String())
		}
		return true, "", err
	}

	return true, strings.TrimSpace(string(output)), nil
}
