package interchaintest_test

import (
	"context"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/relayer/rly"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
)

func TestInterchain_DuplicateChain_CosmosRly(t *testing.T) {
	duplicateChainTest(t, ibc.CosmosRly)
}

func TestInterchain_DuplicateChain_HermesRelayer(t *testing.T) {
	duplicateChainTest(t, ibc.Hermes)
}

func duplicateChainTest(t *testing.T, relayerImpl ibc.RelayerImplementation) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	client, network := interchaintest.DockerSetup(t)

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		// Two otherwise identical chains that only differ by ChainID.
		{Name: "gaia", ChainName: "g1", Version: "v7.0.1"},
		{Name: "gaia", ChainName: "g2", Version: "v7.0.1"},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	gaia0, gaia1 := chains[0], chains[1]

	r := interchaintest.NewBuiltinRelayerFactory(relayerImpl, zaptest.NewLogger(t)).Build(
		t, client, network,
	)

	ic := interchaintest.NewInterchain().
		AddChain(gaia0).
		AddChain(gaia1).
		AddRelayer(r, "r").
		AddLink(interchaintest.InterchainLink{
			Chain1:  gaia0,
			Chain2:  gaia1,
			Relayer: r,
		})

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()
	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:  t.Name(),
		Client:    client,
		NetworkID: network,

		SkipPathCreation: true,
	}))
	_ = ic.Close()
}

func TestInterchain_GetRelayerWallets_CosmosRly(t *testing.T) {
	getRelayerWalletsTest(t, ibc.CosmosRly)
}

func TestInterchain_GetRelayerWallets_HermesRelayer(t *testing.T) {
	getRelayerWalletsTest(t, ibc.Hermes)
}

func getRelayerWalletsTest(t *testing.T, relayerImpl ibc.RelayerImplementation) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	client, network := interchaintest.DockerSetup(t)

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		// Two otherwise identical chains that only differ by ChainID.
		{Name: "gaia", ChainName: "g1", Version: "v7.0.1", ChainConfig: ibc.ChainConfig{ChainID: "cosmoshub-0"}},
		{Name: "gaia", ChainName: "g2", Version: "v7.0.1", ChainConfig: ibc.ChainConfig{ChainID: "cosmoshub-1"}},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	gaia0, gaia1 := chains[0], chains[1]

	r := interchaintest.NewBuiltinRelayerFactory(relayerImpl, zaptest.NewLogger(t)).Build(
		t, client, network,
	)

	ic := interchaintest.NewInterchain().
		AddChain(gaia0).
		AddChain(gaia1).
		AddRelayer(r, "r").
		AddLink(interchaintest.InterchainLink{
			Chain1:  gaia0,
			Chain2:  gaia1,
			Relayer: r,
		})

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()
	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:  t.Name(),
		Client:    client,
		NetworkID: network,

		SkipPathCreation: true,
	}))

	var (
		g1Wallet    ibc.Wallet
		g2Wallet    ibc.Wallet
		walletFound bool
	)

	t.Run("Chain one wallet is returned", func(t *testing.T) {
		g1Wallet, walletFound = r.GetWallet(chains[0].Config().ChainID)
		require.True(t, walletFound)
		require.NotEmpty(t, g1Wallet.Address())
		require.NotEmpty(t, g1Wallet.Mnemonic())
	})

	t.Run("Chain two wallet is returned", func(t *testing.T) {
		g2Wallet, walletFound = r.GetWallet(chains[1].Config().ChainID)
		require.True(t, walletFound)
		require.NotEmpty(t, g2Wallet.Address())
		require.NotEmpty(t, g2Wallet.Mnemonic())
	})

	t.Run("Different wallets are returned", func(t *testing.T) {
		require.NotEqual(t, g1Wallet.Address(), g2Wallet.Address())
		require.NotEqual(t, g1Wallet.Mnemonic(), g2Wallet.Mnemonic())
	})

	t.Run("Wallet for different chain does not exist", func(t *testing.T) {
		_, ok := r.GetWallet("cosmoshub-does-not-exist")
		require.False(t, ok)
	})

	_ = ic.Close()
}

func TestInterchain_CreateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	client, network := interchaintest.DockerSetup(t)

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		// Two otherwise identical chains that only differ by ChainID.
		{Name: "gaia", ChainName: "g1", Version: "v7.0.1", ChainConfig: ibc.ChainConfig{ChainID: "cosmoshub-0"}},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	gaia0 := chains[0]

	ic := interchaintest.NewInterchain().AddChain(gaia0)
	defer ic.Close()

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()
	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:  t.Name(),
		Client:    client,
		NetworkID: network,
	}))

	initBal := math.NewInt(10_000)

	t.Run("with mnemonic", func(t *testing.T) {
		keyName := "mnemonic-user-name"

		registry := codectypes.NewInterfaceRegistry()
		cryptocodec.RegisterInterfaces(registry)
		cdc := codec.NewProtoCodec(registry)

		kr := keyring.NewInMemory(cdc)
		_, mnemonic, err := kr.NewMnemonic(
			keyName,
			keyring.English,
			hd.CreateHDPath(sdk.CoinType, 0, 0).String(),
			"", // Empty passphrase.
			hd.Secp256k1,
		)

		require.NoError(t, err)
		require.NotEmpty(t, mnemonic)

		user, err := interchaintest.GetAndFundTestUserWithMnemonic(ctx, keyName, mnemonic, initBal.Int64(), gaia0)
		require.NoError(t, err)
		require.NoError(t, testutil.WaitForBlocks(ctx, 2, gaia0))
		require.NotEmpty(t, user.Address())
		require.NotEmpty(t, user.KeyName())

		actualBalance, err := gaia0.GetBalance(ctx, user.FormattedAddress(), gaia0.Config().Denom)
		require.NoError(t, err)
		require.True(t, actualBalance.Equal(initBal))
	})

	t.Run("without mnemonic", func(t *testing.T) {
		keyName := "regular-user-name"
		users := interchaintest.GetAndFundTestUsers(t, ctx, keyName, initBal.Int64(), gaia0)
		require.NoError(t, testutil.WaitForBlocks(ctx, 2, gaia0))
		require.Len(t, users, 1)
		require.NotEmpty(t, users[0].Address())
		require.NotEmpty(t, users[0].KeyName())

		actualBalance, err := gaia0.GetBalance(ctx, users[0].FormattedAddress(), gaia0.Config().Denom)
		require.NoError(t, err)
		require.True(t, actualBalance.Equal(initBal))
	})
}

func TestCosmosChain_BroadcastTx_CosmosRly(t *testing.T) {
	broadcastTxCosmosChainTest(t, ibc.CosmosRly)
}

func TestCosmosChain_BroadcastTx_HermesRelayer(t *testing.T) {
	broadcastTxCosmosChainTest(t, ibc.Hermes)
}

func broadcastTxCosmosChainTest(t *testing.T, relayerImpl ibc.RelayerImplementation) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	client, network := interchaintest.DockerSetup(t)

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		// Two otherwise identical chains that only differ by ChainID.
		{Name: "gaia", ChainName: "g1", Version: "v7.0.1", ChainConfig: ibc.ChainConfig{ChainID: "cosmoshub-0"}},
		{Name: "gaia", ChainName: "g2", Version: "v7.0.1", ChainConfig: ibc.ChainConfig{ChainID: "cosmoshub-1"}},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	gaia0, gaia1 := chains[0], chains[1]

	r := interchaintest.NewBuiltinRelayerFactory(relayerImpl, zaptest.NewLogger(t)).Build(
		t, client, network,
	)

	pathName := "p"
	ic := interchaintest.NewInterchain().
		AddChain(gaia0).
		AddChain(gaia1).
		AddRelayer(r, "r").
		AddLink(interchaintest.InterchainLink{
			Chain1:  gaia0,
			Chain2:  gaia1,
			Relayer: r,
			Path:    pathName,
		})

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()
	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:  t.Name(),
		Client:    client,
		NetworkID: network,
	}))

	testUser := interchaintest.GetAndFundTestUsers(t, ctx, "gaia-user-1", 10_000_000, gaia0)[0]

	sendAmount := math.NewInt(10_000)

	t.Run("relayer starts", func(t *testing.T) {
		require.NoError(t, r.StartRelayer(ctx, eRep, pathName))
	})

	t.Run("broadcast success", func(t *testing.T) {
		b := cosmos.NewBroadcaster(t, gaia0.(*cosmos.CosmosChain))
		transferAmount := sdk.Coin{Denom: gaia0.Config().Denom, Amount: sendAmount}
		memo := ""

		msg := transfertypes.NewMsgTransfer(
			"transfer",
			"channel-0",
			transferAmount,
			testUser.FormattedAddress(),
			testUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(gaia1.Config().Bech32Prefix),
			clienttypes.NewHeight(1, 1000),
			0,
			memo,
		)
		resp, err := cosmos.BroadcastTx(ctx, b, testUser.(*cosmos.CosmosWallet), msg)
		require.NoError(t, err)
		assertTransactionIsValid(t, resp)
	})

	t.Run("transfer success", func(t *testing.T) {
		require.NoError(t, testutil.WaitForBlocks(ctx, 5, gaia0, gaia1))

		srcDenomTrace := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom("transfer", "channel-0", gaia0.Config().Denom))
		dstIbcDenom := srcDenomTrace.IBCDenom()

		dstFinalBalance, err := gaia1.GetBalance(ctx, testUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(gaia1.Config().Bech32Prefix), dstIbcDenom)
		require.NoError(t, err, "failed to get balance from dest chain")
		require.True(t, dstFinalBalance.Equal(sendAmount))
	})
}

// An external package that imports interchaintest may not provide a GitSha when they provide a BlockDatabaseFile.
// The GitSha field is documented as optional, so this should succeed.
func TestInterchain_OmitGitSHA(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	client, network := interchaintest.DockerSetup(t)

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{Name: "gaia", Version: "v7.0.1"},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)
	gaia := chains[0]

	ic := interchaintest.NewInterchain().
		AddChain(gaia)

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)
	ctx := context.Background()
	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:  t.Name(),
		Client:    client,
		NetworkID: network,

		SkipPathCreation: true,

		BlockDatabaseFile: ":memory:",
	}))
	_ = ic.Close()
}

func TestInterchain_ConflictRejection(t *testing.T) {
	t.Run("duplicate chain", func(t *testing.T) {
		cf := interchaintest.NewBuiltinChainFactory(zap.NewNop(), []*interchaintest.ChainSpec{
			{Name: "gaia", Version: "v7.0.1", ChainConfig: ibc.ChainConfig{ChainID: "cosmoshub-0"}},
		})

		chains, err := cf.Chains(t.Name())
		require.NoError(t, err)
		chain := chains[0]

		exp := fmt.Sprintf("chain %v was already added", chain)
		require.PanicsWithError(t, exp, func() {
			_ = interchaintest.NewInterchain().AddChain(chain).AddChain(chain)
		})
	})

	t.Run("chain name", func(t *testing.T) {
		cf := interchaintest.NewBuiltinChainFactory(zap.NewNop(), []*interchaintest.ChainSpec{
			// Different ChainID, but explicit ChainName used twice.
			{Name: "gaia", ChainName: "g", Version: "v7.0.1", ChainConfig: ibc.ChainConfig{ChainID: "cosmoshub-0"}},
			{Name: "gaia", ChainName: "g", Version: "v7.0.1", ChainConfig: ibc.ChainConfig{ChainID: "cosmoshub-1"}},
		})

		chains, err := cf.Chains(t.Name())
		require.NoError(t, err)

		require.PanicsWithError(t, "a chain with name g already exists", func() {
			_ = interchaintest.NewInterchain().AddChain(chains[0]).AddChain(chains[1])
		})
	})

	t.Run("chain ID", func(t *testing.T) {
		cf := interchaintest.NewBuiltinChainFactory(zap.NewNop(), []*interchaintest.ChainSpec{
			// Valid ChainName but duplicate ChainID.
			{Name: "gaia", ChainName: "g1", Version: "v7.0.1", ChainConfig: ibc.ChainConfig{ChainID: "cosmoshub-0"}},
			{Name: "gaia", ChainName: "g2", Version: "v7.0.1", ChainConfig: ibc.ChainConfig{ChainID: "cosmoshub-0"}},
		})

		chains, err := cf.Chains(t.Name())
		require.NoError(t, err)

		require.PanicsWithError(t, "a chain with ID cosmoshub-0 already exists", func() {
			_ = interchaintest.NewInterchain().AddChain(chains[0]).AddChain(chains[1])
		})
	})

	t.Run("duplicate relayer", func(t *testing.T) {
		var r rly.CosmosRelayer

		exp := fmt.Sprintf("relayer %v was already added", &r)
		require.PanicsWithError(t, exp, func() {
			_ = interchaintest.NewInterchain().AddRelayer(&r, "r1").AddRelayer(&r, "r2")
		})
	})

	t.Run("relayer name", func(t *testing.T) {
		var r1, r2 rly.CosmosRelayer

		require.PanicsWithError(t, "a relayer with name r already exists", func() {
			_ = interchaintest.NewInterchain().AddRelayer(&r1, "r").AddRelayer(&r2, "r")
		})
	})
}

func TestInterchain_AddNil(t *testing.T) {
	require.PanicsWithError(t, "cannot add nil chain", func() {
		_ = interchaintest.NewInterchain().AddChain(nil)
	})

	require.PanicsWithError(t, "cannot add nil relayer", func() {
		_ = interchaintest.NewInterchain().AddRelayer(nil, "r")
	})
}

func assertTransactionIsValid(t *testing.T, resp sdk.TxResponse) {
	require.NotNil(t, resp)
	require.NotEqual(t, 0, resp.GasUsed)
	require.NotEqual(t, 0, resp.GasWanted)
	require.Equal(t, uint32(0), resp.Code)
	require.NotEmpty(t, resp.Data)
	require.NotEmpty(t, resp.TxHash)
	require.NotEmpty(t, resp.Events)
}
