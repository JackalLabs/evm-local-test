package eth

import (
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

type EthAPI struct {
	client *ethclient.Client

	Retries   int
	RetryWait time.Duration
}

type EthGetProofResponse struct {
	StorageHash  string `json:"storageHash"`
	StorageProof []struct {
		Key   string   `json:"key"`
		Proof []string `json:"proof"`
		Value string   `json:"value"`
	} `json:"storageProof"`
	AccountProof []string `json:"accountProof"`
}

func NewEthAPI(rpc string) (EthAPI, error) {
	ethClient, err := ethclient.Dial(rpc)
	if err != nil {
		return EthAPI{}, err
	}

	return EthAPI{
		client:    ethClient,
		Retries:   6,
		RetryWait: 10 * time.Second,
	}, nil
}
