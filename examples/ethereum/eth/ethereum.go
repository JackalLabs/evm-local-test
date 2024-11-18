package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
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
	//BeaconAPIClient *BeaconAPIClient	NOTE: Eureka used beacon for what?

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
