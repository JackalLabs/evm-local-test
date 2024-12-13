package e2esuite

import (
	"context"

	dockerclient "github.com/docker/docker/client"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ictethereum "github.com/strangelove-ventures/interchaintest/v7/chain/ethereum"
	chainconfig "github.com/strangelove-ventures/interchaintest/v7/examples/ethereum/chainconfig"
	"github.com/strangelove-ventures/interchaintest/v7/examples/ethereum/eth"
	logger "github.com/strangelove-ventures/interchaintest/v7/examples/logger"
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

// SetupSuite sets up the chains, mulberry relayer, user accounts, clients, and connections
func (s *TestSuite) SetupSuite(ctx context.Context) {
	logger.InitLogger()

	icChainSpecs := chainconfig.ChainSpecs
	logger.LogInfo(icChainSpecs)

	// At this step, the ibc team use a case statement to decide whether to boot up a POW or POS Eth chain.
	// We might need to do this in the future.

	s.logger = zaptest.NewLogger(s.T())
	s.dockerClient, s.network = interchaintest.DockerSetup(s.T())

	cf := interchaintest.NewBuiltinChainFactory(s.logger, icChainSpecs)

	chains, err := cf.Chains(s.T().Name())
	s.Require().NoError(err)

	// canine-chain should be at index 1
	s.ChainB = chains[1].(*cosmos.CosmosChain)

	ic := interchaintest.NewInterchain()
	for _, chain := range chains {
		ic = ic.AddChain(chain)
	}

	s.ExecRep = testreporter.NewNopReporter().RelayerExecReporter(s.T())

	// TODO: Run this in a goroutine and wait for it to be ready
	s.Require().NoError(ic.Build(ctx, s.ExecRep, interchaintest.InterchainBuildOptions{
		TestName:         s.T().Name(),
		Client:           s.dockerClient,
		NetworkID:        s.network,
		SkipPathCreation: true,
	}))

	// NOTE: We can map all query request types to their gRPC method paths for cosmos chains?
	// Easier/faster than making function(s) for jackal queries?

	anvil := chains[0].(*ictethereum.EthereumChain)

	faucet, err := crypto.ToECDSA(ethcommon.FromHex(anvilFaucetPrivateKey))
	s.Require().NoError(err)

	s.ChainA, err = eth.NewEthereum(ctx, anvil.GetHostRPCAddress(), faucet)
	s.Require().NoError(err)

	// log
	logger.LogInfo("host rpc address: %s\n", anvil.GetHostRPCAddress())

	// TODO: fund jkl users
	// TODO: make a map of proposal IDs?

}
