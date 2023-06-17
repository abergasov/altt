package rpc

import (
	"altt/internal/entities"
	"container/list"
	"errors"
	"log"
	"sync"

	"go.uber.org/zap"
)

var (
	ErrRPCUnsupportedChain   = errors.New("unsupported chain")
	ErrRPCUninitializedChain = errors.New("uninitialized chain")
)

type Service struct {
	rpcs             map[entities.Chain][]string
	usage            map[entities.Chain]*list.List
	usageMU          sync.Mutex
	configuredChains map[entities.Chain]struct{}
}

var s *Service

func NewService(rpcEndpoints map[string][]string) {
	s = &Service{
		rpcs:             make(map[entities.Chain][]string, len(rpcEndpoints)),
		configuredChains: make(map[entities.Chain]struct{}, len(rpcEndpoints)),
		usage:            make(map[entities.Chain]*list.List),
	}
	for chain, rpcList := range rpcEndpoints {
		c, err := entities.ChainFromString(chain)
		if err != nil {
			log.Fatal("invalid chain", err, zap.String("chain", chain))
		}
		s.rpcs[c] = rpcList
		s.configuredChains[c] = struct{}{}
	}

	for chain, rpcs := range s.rpcs {
		s.usage[chain] = list.New()
		for _, rpc := range rpcs {
			s.usage[chain].PushBack(rpc)
		}
	}
}

func ChainAvailable(chain entities.Chain) bool {
	_, ok := s.configuredChains[chain]
	return ok
}

func GetRPC(chain entities.Chain) (string, error) {
	if s == nil {
		log.Fatal("rpc service not initialized")
	}
	s.usageMU.Lock()
	defer s.usageMU.Unlock()
	if s.usage == nil {
		log.Fatal("rpc service not initialized") // nolint:gocritic
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
