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

	err := cmd.Start()
	if err != nil {
		return "", fmt.Errorf("failed to start command: %s, stderr: %s", err, stderr.String())
	}

	go func() {
		err = cmd.Wait()
		if err != nil {
			fmt.Printf("Command exited with error: %s\n", err)
		}
	}()

	return out.String(), nil
}
