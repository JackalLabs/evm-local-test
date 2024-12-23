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

	// s.Require().True(s.Run("Set up environment", func() {
	// 	err := os.Chdir("../..") // Change directories for what?
	// 	s.Require().NoError(err)
	// }))

	// s.Require().True(s.Run("Deploy ethereum contracts", func() {
	// }))
}

func (s *OutpostTestSuite) TestForge() {
	ctx := context.Background()
	s.SetupForgeSuite(ctx)

	s.Require().True(s.Run("forge", func() {
		fmt.Println("made it to the end")
	}))
	time.Sleep(10 * time.Hour)
}
