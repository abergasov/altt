package balancer

import (
	"altt/internal/logger"
	"altt/internal/service/web3/approver"

	"github.com/golang/groupcache/singleflight"
	"go.uber.org/zap"
)

type Service struct {
	group singleflight.Group
	erc20 *approver.Service
	log   logger.AppLogger
}

func NewService(log logger.AppLogger, erc20 *approver.Service) *Service {
	return &Service{
		log:   log.With(zap.String("service", "balancer")),
		erc20: erc20,
	}
}
