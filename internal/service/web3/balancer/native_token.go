package balancer

import (
	"altt/internal/entities"
	"altt/internal/service/rpc"
	"altt/internal/service/web3"
	"altt/internal/utils"
	"context"
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/common"
)

func (s *Service) GetNativeBalance(ctx context.Context, chain entities.Chain, holder common.Address) (*entities.Balance, error) {
	if !rpc.ChainAvailable(chain) {
		return nil, fmt.Errorf("chain %s is not available", chain.String())
	}
	resp, err := s.group.Do(getNativeKey(chain, holder), func() (interface{}, error) {
		connector, err := web3.GetConnector(chain)
		if err != nil {
			return nil, err
		}
		client, err := connector.GetWeb3()
		if err != nil {
			return nil, err
		}
		return client.BalanceAt(ctx, holder, nil)
	})
	if err != nil {
		s.log.Error("failed to get native balance",
			err,
			zap.String("chain", chain.String()),
			zap.String("address", holder.String()),
		)
		return nil, fmt.Errorf("failed to get native balance")
	}
	return &entities.Balance{
		Chain:           chain,
		ChainName:       chain.String(),
		Token:           entities.MapChainToFuel(chain),
		TokenBalance:    utils.ETHFromWei(resp.(*big.Int)), // use unsafe cast here as we know that it's result of group
		TokenBalanceWei: resp.(*big.Int).String(),
	}, nil
}

func getNativeKey(chain entities.Chain, address common.Address) string {
	return fmt.Sprintf("%s-%s", chain.String(), address.String())
}
