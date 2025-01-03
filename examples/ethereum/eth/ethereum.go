package eth

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/exec"
	"strings"

	"cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
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

func (e Ethereum) ForgeCreate(deployer *ecdsa.PrivateKey, contractName, contractPath string) ([]byte, error) {
	// Prepare the forge create command
	cmd := exec.Command("forge", "create",
		fmt.Sprintf("%s:%s", contractPath, contractName), // Format as "path:ContractName"
		"--rpc-url", e.RPC,
		"--private-key", fmt.Sprintf("0x%s", hex.EncodeToString(deployer.D.Bytes())),
		"-vvvv",
	)

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
		return nil, err
	}

	// Get the output as byte slices
	stdoutBytes := stdoutBuf.Bytes()

	return stdoutBytes, nil
}
