package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/strangelove-ventures/interchaintest/v7/examples/ethereum/eth"
)

func (s *OutpostTestSuite) SetupForgeSuite(ctx context.Context) {
	// Start Anvil node
	anvilArgs := []string{"--port", "8545", "--block-time", "1"}
	output, err := eth.ExecuteCommand("anvil", anvilArgs)
	if err != nil {
		fmt.Printf("Error starting Anvil: %s\n", err)
		return
	}
	fmt.Printf("Anvil Output: %s\n", output)

	// Poll for Anvil readiness
	fmt.Println("Waiting for Anvil to become ready...")
	rpcURL := "http://127.0.0.1:8545"
	if err := waitForRPC(rpcURL, 10*time.Second); err != nil {
		fmt.Printf("Error: Anvil did not become ready in time: %s\n", err)
		return
	}
	fmt.Println("Anvil is ready at", rpcURL)
}

func waitForRPC(rpcURL string, timeout time.Duration) error {
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

func (s *OutpostTestSuite) TestForge() {
	ctx := context.Background()
	s.SetupForgeSuite(ctx)

	s.Require().True(s.Run("forge", func() {
		fmt.Println("made it to the end")
	}))
	time.Sleep(10 * time.Hour) // Placeholder for extended testing
}
