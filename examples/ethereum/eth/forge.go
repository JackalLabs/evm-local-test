package eth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"
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

func WaitForRPC(rpcURL string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for RPC at %s", rpcURL)
		case <-ticker.C:
			// Create a JSON-RPC request
			payload := map[string]interface{}{
				"jsonrpc": "2.0",
				"method":  "eth_blockNumber",
				"params":  []interface{}{},
				"id":      1,
			}
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("failed to marshal JSON payload: %w", err)
			}

			// Send the request
			resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(payloadBytes))
			if err != nil {
				continue // Retry on error
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				return nil // Anvil is ready
			}
		}
	}
}
