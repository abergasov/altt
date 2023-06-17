package balancer

import (
	"altt/internal/entities"
	"altt/internal/service/rpc"
	"altt/internal/service/web3"
	"context"
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/common"
)

func (s *Service) GetKnownTokenBalance(ctx context.Context, token entities.Token, chain entities.Chain, holder common.Address) (*entities.Balance, error) {
	s.metrics.NewTokenBalanceRequest(holder, token)
	if !rpc.ChainAvailable(chain) {
		return nil, fmt.Errorf("chain %s is not available", chain.String())
	}
	resp, err := s.group.Do(getKnownKey(token, chain, holder), func() (interface{}, error) {
		connector, err := web3.GetConnector(chain)
		if err != nil {
			return nil, err
		}
		client, err := connector.GetWeb3()
		if err != nil {
			return nil, fmt.Errorf("unable to get web3 client: %w", err)
		}
		tokenAddress, err := entities.GetTokenAddress(connector.GetChainID(), token)
		if err != nil {
			return nil, fmt.Errorf("unable to get token address: %w", err)
		}
		return s.erc20.GetERC20TokenBalance(ctx, client, tokenAddress, holder)
	})
	if err != nil {
		s.log.Error("failed to get native balance",
			err,
			zap.String("token", string(token)),
			zap.String("chain", chain.String()),
			zap.String("address", holder.String()),
		)
		return nil, fmt.Errorf("failed to get native balance")
	}
	return &entities.Balance{
		Chain:           chain,
		ChainName:       chain.String(),
		Token:           token,
		TokenBalance:    entities.CoinFromWEI(token, resp.(*big.Int)),
		TokenBalanceWei: resp.(*big.Int).String(), // use unsafe cast here as we know that it's result of group
	}, nil
}

func getKnownKey(token entities.Token, chain entities.Chain, address common.Address) string {
	return fmt.Sprintf("%s-%s-%s", string(token), chain.String(), address.String())
}
