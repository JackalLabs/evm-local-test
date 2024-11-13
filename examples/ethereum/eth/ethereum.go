package eth

import (
	"crypto/ecdsa"
	"math/big"
)

// NOTE: This is a 'wrapper' object that works in conjunction with the 'EthereumChain' object
// found in /chain/ethereum/ethereum_chain.go
type Ethereum struct {
	ChainID *big.Int
	RPC     string
	//EthAPI          EthAPI
	//BeaconAPIClient *BeaconAPIClient	NOTE: Eureka used beacon for what?

	Faucet *ecdsa.PrivateKey
}
