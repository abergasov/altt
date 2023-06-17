package main

import (
	"altt/internal/config"
	"altt/internal/logger"
	"altt/internal/routes"
	"altt/internal/service/rpc"
	"altt/internal/service/web3/approver"
	"altt/internal/service/web3/balancer"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

var (
	confFile = flag.String("config", "configs/app_conf.yml", "Configs file path")
	appHash  = os.Getenv("GIT_HASH")
)

func main() {
	flag.Parse()
	appLog, err := logger.NewAppLogger(appHash)
	if err != nil {
		log.Fatalf("unable to create logger: %s", err)
	}
	appLog.Info("app starting", zap.String("conf", *confFile))
	appConf, err := config.InitConf(*confFile)
	if err != nil {
		appLog.Fatal("unable to init config", err, zap.String("config", *confFile))
	}

	appLog.Info("init services")
	rpc.NewService(appConf.ChainRPCs)
	serviceBalancer := balancer.NewService(appLog, approver.InitService(appLog), appConf.DisableMetrics)

	appLog.Info("init http service")
	appHTTPServer := routes.InitAppRouter(appLog, serviceBalancer, fmt.Sprintf(":%d", appConf.AppPort), appConf.DisableMetrics)
	defer func() {
		if err = appHTTPServer.Stop(); err != nil {
			appLog.Fatal("unable to stop http service", err)
		}
	}()
	go func() {
		if err = appHTTPServer.Run(); err != nil {
			appLog.Fatal("unable to start http service", err)
		}
	}()

	// register app shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // This blocks the main thread until an interrupt is received
}
