package utils

import (
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
		return true, "", err
	}

	return true, strings.TrimSpace(string(output)), nil
}
