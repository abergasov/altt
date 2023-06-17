package routes

import (
	"altt/internal/logger"
	"altt/internal/service/web3/balancer"

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
func InitAppRouter(log logger.AppLogger, serviceBalancer *balancer.Service, address string) *Server {
	app := &Server{
		appAddr:         address,
		httpEngine:      fiber.New(fiber.Config{}),
		serviceBalancer: serviceBalancer,
		log:             log.With(zap.String("service", "http")),
	}
	app.httpEngine.Use(recover.New())
	app.initRoutes()
	return app
}

func (s *Server) initRoutes() {
	s.httpEngine.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
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
