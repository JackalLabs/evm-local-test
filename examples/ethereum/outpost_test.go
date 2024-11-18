package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
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
	s.TestSuite.SetupSuite(ctx)

	eth, canined := s.ChainA, s.ChainB
	fmt.Println(eth)
	fmt.Println(canined)

}

func TestWithOutpostTestSuite(t *testing.T) {
	suite.Run(t, new(OutpostTestSuite))
}

func (s *OutpostTestSuite) TestDummy() {
	ctx := context.Background()
	s.SetupSuite(ctx)

	canined := s.ChainB
	fmt.Println(canined)
	time.Sleep(10 * time.Hour)

	s.Require().True(s.Run("dummy", func() {

		fmt.Println("made it here")
		time.Sleep(10 * time.Hour)

	}))
}
