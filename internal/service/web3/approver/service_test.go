package approver_test

import (
	"altt/internal/entities"
	"altt/internal/logger"
	"altt/internal/service/rpc"
	"altt/internal/service/web3"
	"altt/internal/service/web3/approver"
	"altt/internal/utils"
	"context"
	"crypto/ecdsa"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

var (
	sampleContract = "0x8731d54E9D02c286767d56ac03e8037C07e01e98" // in eth mainnet
	targetChain    = entities.ChainEthereum
	sampleRPC      = map[string][]string{
		entities.ChainEthereum.String(): {
			"https://eth.llamarpc.com",
			"https://uk.rpc.blxrbdn.com",
			"https://virginia.rpc.blxrbdn.com",
			"https://rpc.ankr.com/eth",
		},
	}
)

func TestService_ApproveContractUsage(t *testing.T) {
	t.Skip("skip test")
	appLog, err := logger.NewAppLogger("")
	require.NoError(t, err)
	service := approver.InitService(appLog)
	ethClient, privateKey, accAddress := initTest(t)

	tokenAddress, err := entities.GetTokenAddress(targetChain, entities.USDC)
	require.NoError(t, err)

	approveTxHash, err := service.ApproveContractUsageALL(
		ethClient,
		privateKey,
		tokenAddress,
		accAddress,
		common.HexToAddress(sampleContract),
	)
	require.NoError(t, err)
	t.Log("approveTxHash:", approveTxHash)
}

func TestService_GetNativeTokenBalance(t *testing.T) {
	appLog, err := logger.NewAppLogger("")
	require.NoError(t, err)
	service := approver.InitService(appLog)
	ethClient, _, accAddress := initTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	val, err := service.GetNativeTokenBalance(ctx, ethClient, accAddress)
	require.NoError(t, err)
	t.Log("val:", utils.ETHFromWei(val))
}

func TestService_GetERC20TokenBalance(t *testing.T) {
	appLog, err := logger.NewAppLogger("")
	require.NoError(t, err)
	service := approver.InitService(appLog)
	ethClient, _, accAddress := initTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tokenAddress, err := entities.GetTokenAddress(targetChain, entities.USDC)
	require.NoError(t, err)
	val, err := service.GetERC20TokenBalance(ctx, ethClient, tokenAddress, accAddress)
	require.NoError(t, err)
	t.Log("val:", entities.CoinFromWEI(entities.USDC, val))
}

func TestService_GetContractData(t *testing.T) {
	appLog, err := logger.NewAppLogger("")
	require.NoError(t, err)
	service := approver.InitService(appLog)
	ethClient, _, _ := initTest(t)

	require.NoError(t, err)
	tokenAddress, err := entities.GetTokenAddress(targetChain, entities.USDC)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ticker, decimal, err := service.GetContractData(ctx, ethClient, tokenAddress)
	require.NoError(t, err)
	require.Equal(t, string(entities.USDC), ticker)
	require.Equal(t, entities.GetTokenDecimals(entities.USDC), int(decimal))
}

func initTest(t *testing.T) (*ethclient.Client, *ecdsa.PrivateKey, common.Address) {
	rpc.NewService(sampleRPC)
	connector, err := web3.GetConnector(targetChain)
	require.NoError(t, err)
	client, err := connector.GetWeb3()
	require.NoError(t, err)
	t.Cleanup(func() {
		client.Close()
	})

	walletPK := os.Getenv("PRIVATE_KEY")
	if walletPK == "" {
		t.Skip("PRIVATE_KEY is not set")
	}
	privateKey, err := crypto.HexToECDSA(walletPK)
	require.NoError(t, err)

	holder, err := connector.KeyToAddress(privateKey)
	require.NoError(t, err)
	return client, privateKey, holder
}
