package ethereum_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/ethereum"

	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestEthereum(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	client, network := interchaintest.DockerSetup(t)

	// Log location
	f, err := interchaintest.CreateLogFile(fmt.Sprintf("%d.json", time.Now().Unix()))
	require.NoError(t, err)
	// Reporter/logs
	rep := testreporter.NewReporter(f)
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()

	// Get default ethereum chain config for anvil
	anvilConfig := ethereum.DefaultEthereumAnvilChainConfig("ethereum")

	// add --load-state config (this step is not required for tests that don't require an existing state)
	configFileOverrides := make(map[string]any)
	configFileOverrides["--load-state"] = "eigenlayer-deployed-anvil-state.json" // Relative path of state.json
	anvilConfig.ConfigFileOverrides = configFileOverrides

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			ChainName:   "ethereum",
			Name:        "ethereum",
			Version:     "latest",
			ChainConfig: anvilConfig,
		},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	ethereumChain := chains[0].(*ethereum.EthereumChain)

	ic := interchaintest.NewInterchain().
		AddChain(ethereumChain)

	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:  t.Name(),
		Client:    client,
		NetworkID: network,
		// BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation: true, // Skip path creation, so we can have granular control over the process
	}))
	fmt.Println("Interchain built")

	// Sleep for an additional testing
	time.Sleep(10 * time.Hour)

}
