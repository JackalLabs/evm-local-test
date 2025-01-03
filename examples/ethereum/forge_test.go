package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/strangelove-ventures/interchaintest/v7/examples/ethereum/eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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

	// Connect to Anvil RPC
	rpcURL := "http://127.0.0.1:8545"
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// Define accounts and their private keys
	privateKeyA := "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	// Recipient B address
	addressB := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	fmt.Println(addressB)

	// Convert accountA's private key string to *ecdsa.PrivateKey
	privKeyA, err := crypto.HexToECDSA(privateKeyA[2:]) // Remove "0x" prefix
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Get the public address of Account A
	addressA := crypto.PubkeyToAddress(privKeyA.PublicKey)
	fmt.Println(addressA)

	// Check Account A's nonce
	nonce, err := client.PendingNonceAt(context.Background(), addressA)
	if err != nil {
		log.Fatalf("Failed to get nonce for Account A: %v", err)
	}
	fmt.Printf("Account A's nonce is %d\n", nonce)

	// Get chain ID from the client
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	// Prepare the transaction
	amount := new(big.Int).Mul(big.NewInt(35), big.NewInt(1e18)) // 35 ETH in wei
	gasLimit := uint64(21000)                                    // Standard gas limit for ETH transfer
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}

	tx := types.NewTransaction(nonce, addressB, amount, gasLimit, gasPrice, nil)

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privKeyA)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}
	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())

	// Check Account A's nonce
	nonce, err = client.PendingNonceAt(context.Background(), addressA)
	if err != nil {
		log.Fatalf("Failed to get nonce for Account A: %v", err)
	}
	fmt.Printf("Account A's nonce is %d\n", nonce)

	s.Require().True(s.Run("forge", func() {
		fmt.Println("made it to the end")
	}))
	time.Sleep(10 * time.Hour) // Placeholder for extended testing
}
