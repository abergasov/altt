package testhelpers

import (
	"altt/internal/config"
	"altt/internal/entities"
	"altt/internal/logger"
	"altt/internal/service/rpc"
	"altt/internal/service/web3/approver"
	"altt/internal/service/web3/balancer"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestContainer struct {
	Log  logger.AppLogger
	Conf *config.AppConfig

	ServiceBalancer *balancer.Service
}

func GetClean(t *testing.T) *TestContainer {
	conf := getTestConfig()
	rpc.NewService(conf.ChainRPCs)

	appLog, err := logger.NewAppLogger("test")
	require.NoError(t, err)

	serviceBalancer := balancer.NewService(appLog, approver.InitService(appLog), conf.DisableMetrics)

	return &TestContainer{
		Log:             appLog,
		Conf:            conf,
		ServiceBalancer: serviceBalancer,
	}
}

func getTestConfig() *config.AppConfig {
	return &config.AppConfig{
		DisableMetrics: true,
		AppPort:        0,
		ChainRPCs: map[string][]string{
			entities.ChainEthereum.String(): {
				"https://eth.llamarpc.com",
				"https://uk.rpc.blxrbdn.com",
				"https://virginia.rpc.blxrbdn.com",
				"https://rpc.ankr.com/eth",
			},
		},
	}
}
