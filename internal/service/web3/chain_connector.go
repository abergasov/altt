package web3

import (
	"altt/internal/entities"
	"altt/internal/service/rpc"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

var (
	MapChainConnector = map[entities.Chain]ChainConnector{
		entities.ChainCelo: {
			Explorer: "https://celoscan.io/",
			ChainId:  entities.ChainCelo,
		},
		entities.ChainGnosis: {
			Explorer: "https://gnosisscan.io/",
			ChainId:  entities.ChainGnosis,
		},
		entities.ChainHarmony: {
			Explorer: "https://explorer.harmony.one",
			ChainId:  entities.ChainHarmony,
		},
		entities.ChainEthereum: {
			Explorer: "https://etherscan.io",
			ChainId:  entities.ChainEthereum,
		},
		entities.ChainPolygon: {
			Explorer: "https://polygonscan.com",
			ChainId:  entities.ChainPolygon,
		},
		entities.ChainArbitrum: {
			Explorer: "https://arbiscan.io",
			ChainId:  entities.ChainArbitrum,
		},
		entities.ChainOptimism: {
			Explorer: "https://optimistic.etherscan.io",
			ChainId:  entities.ChainOptimism,
		},
		entities.ChainFantom: {
			Explorer: "https://ftmscan.com",
			ChainId:  entities.ChainFantom,
		},
		entities.ChainAvalanche: {
			Explorer: "https://snowtrace.io",
			ChainId:  entities.ChainAvalanche,
		},
		entities.ChainBNB: {
			Explorer: "https://bscscan.com",
			ChainId:  entities.ChainBNB,
		},
		entities.ChainCoreDAO: {
			Explorer: "https://scan.coredao.org/",
			ChainId:  entities.ChainCoreDAO,
		},
	}
)

type ChainConnector struct {
	Explorer string
	ChainId  entities.Chain
}

func (s *ChainConnector) GetWeb3() (*ethclient.Client, error) {
	rpcURL, err := rpc.GetRPC(s.ChainId)
	if err != nil {
		return nil, fmt.Errorf("get rpc url: %w", err)
	}
	return ethclient.Dial(rpcURL)
}

func (c *ChainConnector) GetChainId() entities.Chain {
	return c.ChainId
}

func (s *ChainConnector) GetChainIDBI() *big.Int {
	return big.NewInt(int64(s.ChainId))
}

func (c *ChainConnector) TxHash(hash string) string {
	return c.Explorer + "/tx/" + hash
}
