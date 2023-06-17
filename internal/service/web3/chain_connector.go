package web3

import (
	"altt/internal/entities"
	"altt/internal/service/rpc"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	MapChainConnector = map[entities.Chain]ChainConnector{
		entities.ChainCelo: {
			explorer: "https://celoscan.io/",
			chainID:  entities.ChainCelo,
		},
		entities.ChainGnosis: {
			explorer: "https://gnosisscan.io/",
			chainID:  entities.ChainGnosis,
		},
		entities.ChainHarmony: {
			explorer: "https://explorer.harmony.one",
			chainID:  entities.ChainHarmony,
		},
		entities.ChainEthereum: {
			explorer: "https://etherscan.io",
			chainID:  entities.ChainEthereum,
		},
		entities.ChainPolygon: {
			explorer: "https://polygonscan.com",
			chainID:  entities.ChainPolygon,
		},
		entities.ChainArbitrum: {
			explorer: "https://arbiscan.io",
			chainID:  entities.ChainArbitrum,
		},
		entities.ChainOptimism: {
			explorer: "https://optimistic.etherscan.io",
			chainID:  entities.ChainOptimism,
		},
		entities.ChainFantom: {
			explorer: "https://ftmscan.com",
			chainID:  entities.ChainFantom,
		},
		entities.ChainAvalanche: {
			explorer: "https://snowtrace.io",
			chainID:  entities.ChainAvalanche,
		},
		entities.ChainBNB: {
			explorer: "https://bscscan.com",
			chainID:  entities.ChainBNB,
		},
		entities.ChainCoreDAO: {
			explorer: "https://scan.coredao.org/",
			chainID:  entities.ChainCoreDAO,
		},
	}
)

type ChainConnector struct {
	explorer string
	chainID  entities.Chain
}

func GetConnector(chain entities.Chain) (*ChainConnector, error) {
	conn, ok := MapChainConnector[chain]
	if !ok {
		return nil, fmt.Errorf("chain %s not supported", chain)
	}
	return &conn, nil
}

func (c *ChainConnector) GetWeb3() (*ethclient.Client, error) {
	rpcURL, err := rpc.GetRPC(c.chainID)
	if err != nil {
		return nil, fmt.Errorf("get rpc url: %w", err)
	}
	return ethclient.Dial(rpcURL)
}

func (c *ChainConnector) GetChainID() entities.Chain {
	return c.chainID
}

func (c *ChainConnector) GetChainIDBI() *big.Int {
	return big.NewInt(int64(c.chainID))
}

func (c *ChainConnector) TxHash(hash string) string {
	return c.explorer + "/tx/" + hash
}

func (c *ChainConnector) KeyToAddress(key *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}
