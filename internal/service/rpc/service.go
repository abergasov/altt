package rpc

import (
	"altt/internal/entities"
	"altt/internal/logger"
	"container/list"
	"errors"
	"go.uber.org/zap"
	"log"
	"sync"
)

var (
	ErrRPCUnsupportedChain   = errors.New("unsupported chain")
	ErrRPCUninitializedChain = errors.New("uninitialized chain")
)

type Service struct {
	log     logger.AppLogger
	rpcs    map[entities.Chain][]string
	usage   map[entities.Chain]*list.List
	usageMU sync.Mutex
}

var s *Service

func NewService(log logger.AppLogger, rpcEndpoints map[string][]string) {
	preparedRPCs := make(map[entities.Chain][]string, len(rpcEndpoints))
	for chain, rpcList := range rpcEndpoints {
		c, err := entities.ChainFromString(chain)
		if err != nil {
			log.Fatal("invalid chain", err, zap.String("chain", chain))
		}
		preparedRPCs[c] = rpcList
	}
	s = &Service{
		log:   log.With(zap.String("service", "rpc")),
		rpcs:  preparedRPCs,
		usage: make(map[entities.Chain]*list.List),
	}
	for chain, rpcs := range s.rpcs {
		s.usage[chain] = list.New()
		for _, rpc := range rpcs {
			s.usage[chain].PushBack(rpc)
		}
	}
}

func GetRPC(chain entities.Chain) (string, error) {
	s.usageMU.Lock()
	defer s.usageMU.Unlock()
	if s.usage == nil {
		log.Fatal("rpc service not initialized")
	}
	if _, ok := s.usage[chain]; !ok {
		return "", ErrRPCUnsupportedChain
	}
	if s.usage[chain].Len() == 0 {
		return "", ErrRPCUninitializedChain
	}
	e := s.usage[chain].Front()
	s.usage[chain].MoveToBack(e)
	return e.Value.(string), nil
}
