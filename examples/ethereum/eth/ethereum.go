package eth

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"os/exec"
	"strings"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/strangelove-ventures/interchaintest/v7/examples/ethereum/testvalues"
)

// NOTE: This is a 'wrapper' object that works in conjunction with the 'EthereumChain' object
// found in /chain/ethereum/ethereum_chain.go
type Ethereum struct {
	ChainID *big.Int
	RPC     string
	EthAPI  EthAPI
	// BeaconAPIClient *BeaconAPIClient	NOTE: Eureka used beacon for what?

	Faucet *ecdsa.PrivateKey
}

func NewEthereum(ctx context.Context, rpc string, faucet *ecdsa.PrivateKey) (Ethereum, error) {
	ethClient, err := ethclient.Dial(rpc)
	if err != nil {
		return Ethereum{}, err
	}
	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return Ethereum{}, err
	}

	ethAPI, err := NewEthAPI(rpc)
	if err != nil {
		return Ethereum{}, err
	}

	return Ethereum{
		ChainID: chainID,
		RPC:     rpc,
		EthAPI:  ethAPI,
		Faucet:  faucet,
	}, nil
}

func (e Ethereum) CreateAndFundUser() (*ecdsa.PrivateKey, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	if err := e.FundUser(address, testvalues.StartingEthBalance); err != nil {
		return nil, err
	}

	return key, nil
}

func (e Ethereum) FundUser(address string, amount math.Int) error {
	return e.SendEth(e.Faucet, address, amount)
}

func (e Ethereum) SendEth(key *ecdsa.PrivateKey, toAddress string, amount math.Int) error {
	cmd := exec.Command(
		"cast",
		"send",
		toAddress,
		"--value", amount.String(),
		"--private-key", fmt.Sprintf("0x%s", ethcommon.Bytes2Hex(key.D.Bytes())),
		"--rpc-url", e.RPC,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send eth with %s: %w", strings.Join(cmd.Args, " "), err)
	}

	return nil
}

func (e Ethereum) ForgeScript(deployer *ecdsa.PrivateKey, solidityContract string) ([]byte, error) {
	cmd := exec.Command("forge", "script", "--rpc-url", e.RPC, "--broadcast", "--non-interactive", "-vvvv", solidityContract)

	faucetAddress := crypto.PubkeyToAddress(e.Faucet.PublicKey)
	extraEnv := []string{
		fmt.Sprintf("%s=%s", testvalues.EnvKeyE2EFaucetAddress, faucetAddress.Hex()),
		fmt.Sprintf("PRIVATE_KEY=0x%s", hex.EncodeToString(deployer.D.Bytes())),
	}

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, extraEnv...)

	var stdoutBuf bytes.Buffer

	// Create a MultiWriter to write to both os.Stdout and the buffer
	multiWriter := io.MultiWriter(os.Stdout, &stdoutBuf)

	// Set the command's stdout to the MultiWriter
	cmd.Stdout = multiWriter
	cmd.Stderr = os.Stderr
	fmt.Println("The args are", cmd.Args)

	// Run the command
	if err := cmd.Run(); err != nil {
		fmt.Println("Error start command", cmd.Args, err)
		return nil, err
	}

	// Get the output as byte slices
	stdoutBytes := stdoutBuf.Bytes()

	return stdoutBytes, nil
}

func (e Ethereum) ForgeCreate(deployer *ecdsa.PrivateKey, contractName, contractPath string, relays []string, priceFeed string) (string, error) {
	// Prepare the forge create command
	relaysArg := fmt.Sprintf("[%s]", strings.Join(relays, ",")) // Format array as [address1,address2,...]

	cmdArgs := []string{
		"create",
		"--broadcast",
		fmt.Sprintf("%s:%s", contractPath, contractName), // Format as "path:ContractName"
		"--rpc-url", e.RPC,
		"--private-key", fmt.Sprintf("0x%s", hex.EncodeToString(deployer.D.Bytes())),
		"--remappings", "@openzeppelin/contracts-upgradeable/=forge/lib/openzeppelin-contracts-upgradeable/contracts/",
		"--remappings", "@openzeppelin/contracts/=forge/lib/openzeppelin-contracts-upgradeable/lib/openzeppelin-contracts/contracts/",
		"--remappings", "@chainlink/interfaces/=forge/lib/foundry-chainlink-toolkit/src/interfaces/",
		"--gas-price", "20000000000",
		"--constructor-args", relaysArg, priceFeed,
	}

	cmd := exec.Command("forge", cmdArgs...)

	// Inherit the parent process environment
	cmd.Env = os.Environ()

	var stdoutBuf bytes.Buffer

	// Create a MultiWriter to write to both os.Stdout and the buffer
	multiWriter := io.MultiWriter(os.Stdout, &stdoutBuf)

	// Set the command's stdout and stderr to MultiWriter
	cmd.Stdout = multiWriter
	cmd.Stderr = os.Stderr

	// Debugging: Print the command arguments
	fmt.Println("The args are", cmd.Args)

	// Run the command
	if err := cmd.Run(); err != nil {
		fmt.Println("Error executing command:", cmd.Args, err)
		return "", err
	}

	// Parse the output to find the deployed contract address
	output := stdoutBuf.String()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Deployed to:") {
			// Extract the address after "Deployed to:"
			parts := strings.Fields(line)
			if len(parts) > 2 {
				return parts[2], nil // Return the contract address
			}
		}
	}

	// If no address is found, return an error
	return "", fmt.Errorf("could not find deployed contract address in output")
}

func (e Ethereum) CastSend(contractAddress, functionSig string, args []string, rpcURL, privateKey string, value *big.Int) (string, error) {
	// Prepare the `cast send` command
	cmdArgs := []string{"send", contractAddress, functionSig}
	cmdArgs = append(cmdArgs, args...) // Append function arguments
	cmdArgs = append(cmdArgs, "--rpc-url", rpcURL, "--private-key", privateKey)
	if value != nil {
		cmdArgs = append(cmdArgs, "--value", value.String())
	}
	cmdArgs = append(cmdArgs, "--gas-price", "2000000000000000")
	cmdArgs = append(cmdArgs, "--gas-limit", "1000000")

	cmd := exec.Command("cast", cmdArgs...)

	// Log the full command being executed
	fmt.Printf("Executing command: cast %s\n", strings.Join(cmdArgs, " "))

	// Capture output for debugging
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Run the command
	err := cmd.Run()

	// Extract the transaction hash if present
	output := stdoutBuf.String()
	var txHash string
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "transactionHash") {
			txHash = strings.TrimSpace(strings.Fields(line)[1])
			break
		}
	}

	// Print the transaction hash
	fmt.Printf("Transaction hash: %s\n", txHash)

	// Check for errors
	if err != nil {
		fmt.Printf("Error executing cast send: %s\nStdout: %s\nStderr: %s\n", err, output, stderrBuf.String())
		return txHash, err
	}

	// Print successful execution
	fmt.Printf("Successfully called `%s` on contract %s with args %v\nOutput: %s\n",
		functionSig, contractAddress, args, output)
	return txHash, nil
}

// CastCall uses `cast call` to interact with a view function of any Ethereum contract.
// It decodes the output into a human-readable format for standard return types.
func CastCall(contractAddress, functionSig string, rpcURL string, args []string) (string, error) {
	// Prepare the `cast call` command
	cmdArgs := []string{"call", contractAddress, functionSig}
	cmdArgs = append(cmdArgs, args...) // Append function arguments if needed
	cmdArgs = append(cmdArgs, "--rpc-url", rpcURL)

	cmd := exec.Command("cast", cmdArgs...)

	// Capture output for debugging
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Run the command
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error executing cast call: %s\nStdout: %s\nStderr: %s\n", err, stdoutBuf.String(), stderrBuf.String())
		return "", err
	}

	// Process the raw output
	rawOutput := strings.TrimSpace(stdoutBuf.String())

	// Decode the output
	decodedValue, err := decodeHexOutput(rawOutput)
	if err != nil {
		fmt.Printf("Error decoding output: %v\n", err)
		return rawOutput, nil // Return raw output as fallback
	}

	fmt.Printf("Successfully called `%s` on contract %s\nDecoded Output: %s\n", functionSig, contractAddress, decodedValue)
	return decodedValue, nil
}

// decodeHexOutput decodes an ABI-encoded hex string into a human-readable number (for uint256 results).
func decodeHexOutput(hexOutput string) (string, error) {
	// Remove "0x" prefix if present
	if strings.HasPrefix(hexOutput, "0x") {
		hexOutput = hexOutput[2:]
	}

	// Decode the hex string into bytes
	bytes, err := hex.DecodeString(hexOutput)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex string: %w", err)
	}

	// Convert the bytes to a big.Int
	result := new(big.Int).SetBytes(bytes)
	return result.String(), nil
}

func ListenToLogs(client *ethclient.Client, contractAddress common.Address) {
	// Compute the signature hash of the `Debug` event
	eventSignature := []byte("Debug(string)")
	eventSignatureHash := crypto.Keccak256Hash(eventSignature)

	// Set up the filter query
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics: [][]common.Hash{
			{eventSignatureHash}, // Match only `Debug` events
		},
	}

	// Create a channel for logs
	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}

	fmt.Println("Listening for Debug events...")

	// Process logs
	for {
		select {
		case err := <-sub.Err():
			if strings.Contains(err.Error(), "unexpected EOF") {
				log.Printf("Websocket error: %v", err) // triggered on ctrl-c
			} else {
				log.Fatalf("Subscription error: %v", err)
			}
		case vLog := <-logs:
			// Decode the `Debug` event
			if len(vLog.Data) > 0 {
				message := string(vLog.Data)
				fmt.Printf("Debug Event - Message: %s\n", message)
			} else {
				fmt.Printf("Received Debug Event - Raw Log: %+v\n", vLog)
			}
		}
	}
}
