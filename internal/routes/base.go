package routes

import (
	"altt/internal/logger"
	"altt/internal/service/web3/balancer"

	fiberprometheus "github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

type Server struct {
	appAddr         string
	log             logger.AppLogger
	serviceBalancer *balancer.Service
	httpEngine      *fiber.App
}

// InitAppRouter initializes the HTTP Server.
func InitAppRouter(log logger.AppLogger, serviceBalancer *balancer.Service, address string, disableMetrics bool) *Server {
	app := &Server{
		appAddr:         address,
		httpEngine:      fiber.New(fiber.Config{}),
		serviceBalancer: serviceBalancer,
		log:             log.With(zap.String("service", "http")),
	}
	app.httpEngine.Use(recover.New())
	app.initRoutes(disableMetrics)
	return app
}

func (s *Server) initRoutes(disableMetrics bool) {
	if !disableMetrics {
		prometheus := fiberprometheus.New("balancer-app")
		prometheus.RegisterAt(s.httpEngine, "/metrics")
		s.httpEngine.Use(prometheus.Middleware)
	}
	s.httpEngine.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.SendString("ok") // no extra deps, simply return 200
	})
	s.httpEngine.Get("/ready", func(ctx *fiber.Ctx) error {
		return ctx.SendString("ok")
	})
	s.httpEngine.Get("/:chain/balance/:address", s.getNativeBalance)
	s.httpEngine.Get("/:chain/:token/balance/:address", s.getKnownTokenBalance)
}

// Run starts the HTTP Server.
func (s *Server) Run() error {
	s.log.Info("Starting HTTP server", zap.String("port", s.appAddr))
	return s.httpEngine.Listen(s.appAddr)
}

func (s *Server) Stop() error {
	return s.httpEngine.Shutdown()
}
