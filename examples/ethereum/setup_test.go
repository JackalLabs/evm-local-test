package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/examples/ethereum/chainconfig"
	"github.com/strangelove-ventures/interchaintest/v7/examples/ethereum/e2esuite"
	"github.com/strangelove-ventures/interchaintest/v7/examples/ethereum/eth"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"go.uber.org/zap/zaptest"
)

var canineRPCAddress string

func (s *OutpostTestSuite) SetupJackalEVMBridgeSuite(ctx context.Context) {
	// Start Anvil node
	anvilArgs := []string{"--port", "8545", "--block-time", "1"}
	// easiest way to install anvil is foundryup --install stable
	// you can modify the code to use docker container with --network host
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

	// Start canine-chain with 3 validators
	icChainSpecs := chainconfig.ChainSpecs

	s.TestSuite.Logger = zaptest.NewLogger(s.T())
	s.TestSuite.DockerClient, s.Network = interchaintest.DockerSetup(s.T())

	cf := interchaintest.NewBuiltinChainFactory(s.Logger, icChainSpecs)

	chains, err := cf.Chains(s.T().Name())
	s.Require().NoError(err)

	// canine-chain should be at index 0

	ic := interchaintest.NewInterchain()
	for _, chain := range chains {
		ic = ic.AddChain(chain)
	}

	s.TestSuite.ExecRep = testreporter.NewNopReporter().RelayerExecReporter(s.T())

	// TODO: Run this in a goroutine and wait for it to be ready
	s.Require().NoError(ic.Build(ctx, s.ExecRep, interchaintest.InterchainBuildOptions{
		TestName:         s.T().Name(),
		Client:           s.DockerClient,
		NetworkID:        s.Network,
		SkipPathCreation: true,
	}))

	canine := chains[0].(*cosmos.CosmosChain)
	canineRPC := canine.GetRPCAddress()
	canineRPCAddress = canineRPC
	log.Printf("canine-chain rpc is: %s", canineRPCAddress)
	canineHostRPC := canine.GetHostRPCAddress()
	log.Printf("canine-chain host rpc is: %s", canineHostRPC)

	// NOTE: I think Mulberry should be able to listen to canine-chain using '127.0.0.1' now
	// TODO: change it back to local host then
	updatedCanineHostRPC := strings.Replace(canineHostRPC, "127.0.0.1", "host.docker.internal", 1)
	log.Printf("updatedCanineHostRPC is: %s", updatedCanineHostRPC)

	// returned canine-chain rpc is: http://puppy-1-fn-0-TestWithOutpostTestSuite_TestJackalEVMBridge:26657
	// and canine-chain host rpc is: http://127.0.0.1:59026

	// Mulberry just has to ping it using , e.g. http://host.docker.internal:59026 -- recreate this with each run
	// So we should boot canine-chain before mulberry

	// setup Mulberry, pull image
	var image string
	switch runtime.GOARCH {
	case "arm64":
		image = "biphan4/mulberry:0.0.9"
	case "amd64":
		image = "anthonyjackallabs/mulberry"
	default:
		log.Fatalf("unsupported architecture %s", runtime.GOARCH)
	}

	if err := e2esuite.PullMulberryImage(image); err != nil {
		log.Fatalf("Error pulling Docker image: %v", err)
	}

	// Get the absolute path of the local config file
	localConfigPath, err := filepath.Abs("e2esuite/mulberry_config.yaml")
	if err != nil {
		log.Fatalf("failed to resolve config path: %v", err)
	}

	// Run the container, stream logs
	containerID, err := e2esuite.RunContainerWithConfig(image, "mulberry", localConfigPath)
	if err != nil {
		log.Fatalf("Error running container: %v", err)
	}

	logFile, err := os.Create("mulberry_logs.txt")
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	go func() {
		err := e2esuite.StreamContainerLogsToFile(containerID, logFile)
		if err != nil {
			log.Printf("Failed to stream Mulberry logs to file: %v", err)
		}
	}()

	// Give mulberry a wallet
	addressCommand := []string{"sh", "-c", "mulberry wallet address >> /proc/1/fd/1 2>> /proc/1/fd/2"}
	if err := e2esuite.ExecCommandInContainer(containerID, addressCommand); err != nil {
		log.Fatalf("Error creating wallet address in container: %v", err)
	}

	// Update the YAML file to connect with anvil
	rpcAddress := "http://127.0.0.1:8545"
	wsAddress := "ws://host.docker.internal:8545"
	if err := e2esuite.UpdateMulberryConfigRPC(localConfigPath, "Ethereum Sepolia", rpcAddress, wsAddress); err != nil {
		log.Fatalf("Failed to update mulberry config: %v", err)
	}

	log.Printf("Updated mulberry config with WS address: %s\n", wsAddress)

	// TODO: we can put the bindings contract address here?
	// Update the YAML file to connect with canine-chain
	if err := e2esuite.UpdateMulberryJackalConfigRPC(localConfigPath, updatedCanineHostRPC); err != nil {
		log.Fatalf("Failed to update mulberry's jackal config: %v", err)
	}

	// Start Mulberry
	// NOTE: get logs some other way, streaming the output of 'start' is blocking the rest of the code
	startCommand := []string{"sh", "-c", "mulberry start >> /proc/1/fd/1 2>> /proc/1/fd/2"}
	if err := e2esuite.ExecCommandInContainer(containerID, startCommand); err != nil {
		log.Fatalf("Error starting mulberry in container: %v", err)
	}

}
