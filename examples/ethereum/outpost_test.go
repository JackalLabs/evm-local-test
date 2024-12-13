package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/strangelove-ventures/interchaintest/v7/examples/ethereum/e2esuite"
	"github.com/strangelove-ventures/interchaintest/v7/examples/ethereum/eth"
	"github.com/stretchr/testify/suite"
)

type OutpostTestSuite struct {
	e2esuite.TestSuite

	// Whether to generate fixtures for the solidity tests
	generateFixtures bool

	// The private key of a test account
	key *ecdsa.PrivateKey
	// The private key of the faucet account of interchaintest
	deployer *ecdsa.PrivateKey

	contractAddresses eth.DeployedContracts
}

func (s *OutpostTestSuite) SetupSuite(ctx context.Context) {

	// using local image for now
	image := "biphan4/mulberry:0.0.6"
	if err := e2esuite.PullMulberryImage(image); err != nil {
		log.Fatalf("Error pulling Docker image: %v", err)
	}

	containerName := "mulberry_test_container"

	// Run the container
	containerID, err := e2esuite.RunContainer(image, containerName)
	if err != nil {
		log.Fatalf("Error running container: %v", err)
	}

	log.Printf("Container is running with ID: %s\n", containerID)

	// Execute a command inside the container
	command := []string{"sh", "-c", "mulberry wallet address >> /proc/1/fd/1 2>> /proc/1/fd/2"}
	if err := e2esuite.ExecCommandInContainer(containerID, command); err != nil {
		log.Fatalf("Error executing command in container: %v", err)
	}

	fmt.Println("Mulberry running?")
	time.Sleep(10 * time.Hour)

	s.TestSuite.SetupSuite(ctx)

	eth, canined := s.ChainA, s.ChainB
	fmt.Println(eth)
	fmt.Println(canined)

	s.Require().True(s.Run("Set up environment", func() {
		err := os.Chdir("../..") // Change directories for what?
		s.Require().NoError(err)

		s.key, err = eth.CreateAndFundUser()
		s.Require().NoError(err)

		operatorKey, err := eth.CreateAndFundUser()
		fmt.Println(operatorKey)
		s.Require().NoError(err)

		s.deployer, err = eth.CreateAndFundUser()
		s.Require().NoError(err)

	}))

	s.Require().True(s.Run("Deploy ethereum contracts", func() {
		// seems the operator key is for supporting proofs
		// we're not running proofs atm

		var (
			stdout []byte
			err    error
		)

		// note: can't just pick a name--need actual name of contract. This is case sensitive

		/* NOTE:
		We ran the command:
		forge script --rpc-url http://127.0.0.1:52078 --broadcast --non-interactive
		-vvvv /Users/biphan/jackal/ict-evm/examples/ethereum/scripts/SimpleStorage.s.sol:SimpleStorage

		in our local terminal and it worked
		This means the 'ForgeScript' function is actually targeting our local file system,
		which means creating a mount bind between local scripts directory and the container was pointless?
		*/

		dir, _ := os.Getwd() // note: returns the root of this repository: ict-evm/
		pathOfScripts := filepath.Join(dir, "examples/ethereum/scripts/SimpleStorage.s.sol:SimpleStorage")

		stdout, err = eth.ForgeScript(s.deployer, pathOfScripts)
		fmt.Println(stdout)
		fmt.Println(err)
		fmt.Println("****deployment complete****")

	}))
}

func TestWithOutpostTestSuite(t *testing.T) {
	suite.Run(t, new(OutpostTestSuite))
}

func (s *OutpostTestSuite) TestDummy() {
	ctx := context.Background()
	s.SetupSuite(ctx)

	canined := s.ChainB
	fmt.Println(canined)

	s.Require().True(s.Run("dummy", func() {

		fmt.Println("made it here")
		time.Sleep(10 * time.Hour)

	}))
}
