package metrics

import (
	"altt/internal/entities"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "balancer_proxy"
)

type Service struct {
	disableMetrics      bool
	uniqueAddresses     *prometheus.GaugeVec
	uniqueTokens        *prometheus.GaugeVec
	totalNativeRequests prometheus.Counter
	totalTokenRequests  prometheus.Counter

	reg prometheus.Registerer
}

func IniMetrics(disableMetrics bool) *Service {
	srv := &Service{
		disableMetrics: disableMetrics,
		reg:            prometheus.DefaultRegisterer,
	}
	if !disableMetrics {
		srv.totalNativeRequests = srv.registerCounter("total_native_requests", "Total number of requests for native token")
		srv.totalTokenRequests = srv.registerCounter("total_token_requests", "Total number of requests for custom token")
		srv.uniqueAddresses = srv.registerGauge("unique_addresses", "Total number of unique addresses", []string{"address"})
		srv.uniqueTokens = srv.registerGauge("unique_tokens", "Total number of unique tokens", []string{"token"})
	}
	return srv
}

func (s *Service) registerGauge(name, help string, labels []string) *prometheus.GaugeVec {
	return promauto.With(s.reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      name,
		Help:      help,
	}, labels)
}

func (s *Service) registerCounter(name, help string) prometheus.Counter {
	return promauto.With(s.reg).NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      name,
		Help:      help,
	})
}

func (s *Service) NewNativeBalanceRequest(address common.Address) {
	if s.disableMetrics {
		return
	}
	s.totalNativeRequests.Inc()
	s.uniqueAddresses.WithLabelValues(address.String()).Inc()
}

func (s *Service) NewTokenBalanceRequest(address common.Address, token entities.Token) {
	if s.disableMetrics {
		return
	}
	s.totalTokenRequests.Inc()
	s.uniqueAddresses.WithLabelValues(address.String()).Inc()
	s.uniqueTokens.WithLabelValues(string(token)).Inc()
}
