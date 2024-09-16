package report

import "os"

const machineIdPath = "/etc/machine-id"

func GetHostID() (string, error) {
	machineId, err := os.ReadFile(machineIdPath)
	if err != nil {
		return "", err
	}

	return string(machineId), nil
}
