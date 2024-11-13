package e2esuite

import (
	dockerclient "github.com/docker/docker/client"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"

	"github.com/strangelove-ventures/interchaintest/v7/examples/ethereum/eth"
)

// Is this a new one or the one that already exists in eigenlayer-deployed-anvil-state.json
const anvilFaucetPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

type TestSuite struct {
	suite.Suite

	ChainA         eth.Ethereum
	ethTestnetType string
	ChainB         *cosmos.CosmosChain
	UserB          ibc.Wallet // At some point we will introduce the jackal user
	dockerClient   *dockerclient.Client
	network        string
	logger         *zap.Logger
	ExecRep        *testreporter.RelayerExecReporter

	// Don't need light clients for now. Only concerned about deploying outpost and
	// emitting events

}
