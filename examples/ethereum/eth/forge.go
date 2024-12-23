package eth

import (
	"bytes"
	"fmt"
	"os/exec"
)

func ExecuteCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Start() // Start the command without waiting for it to finish
	if err != nil {
		return "", fmt.Errorf("failed to start command: %s", err)
	}

	// Optionally, you can wait for a short time to ensure it's started correctly
	// Or monitor its output asynchronously using goroutines
	return out.String(), nil
}
