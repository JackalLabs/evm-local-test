package main

import (
	"context"
	"fmt"
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
	if err := eth.WaitForRPC(rpcURL, 10*time.Second); err != nil {
		fmt.Printf("Error: Anvil did not become ready in time: %s\n", err)
		return
	}
	fmt.Println("Anvil is ready at", rpcURL)
}

func (s *OutpostTestSuite) TestForge() {
	ctx := context.Background()
	s.SetupForgeSuite(ctx)

	s.Require().True(s.Run("forge", func() {
		fmt.Println("made it to the end")
	}))
	time.Sleep(10 * time.Hour) // Placeholder for extended testing
}
