package main

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-health/internal/config"
	"github.com/RacoonMediaServer/rms-health/internal/monitor"
	"github.com/RacoonMediaServer/rms-health/internal/service"
	"github.com/RacoonMediaServer/rms-packages/pkg/pubsub"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"time"

	// Plugins
	_ "github.com/go-micro/plugins/v4/registry/etcd"
)

var Version = "v0.0.0"

const serviceName = "rms-health"

func main() {
	logger.Infof("%s %s", serviceName, Version)
	defer logger.Info("DONE.")

	useDebug := false

	microService := micro.NewService(
		micro.Name(serviceName),
		micro.Version(Version),
		micro.Flags(
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"debug"},
				Usage:       "debug log level",
				Value:       false,
				Destination: &useDebug,
			},
		),
	)

	microService.Init(
		micro.Action(func(context *cli.Context) error {
			configFile := fmt.Sprintf("/etc/rms/%s.json", serviceName)
			if context.IsSet("config") {
				configFile = context.String("config")
			}
			return config.Load(configFile)
		}),
	)

	if useDebug {
		_ = logger.Init(logger.WithLevel(logger.DebugLevel))
	}

	cfg := config.Config()
	healthService := service.New(pubsub.NewPublisher(microService))

	mon := monitor.Monitor{
		Factory:       servicemgr.NewServiceFactory(microService),
		CheckInterval: time.Duration(cfg.CheckIntervalMin) * time.Second,
		ReportChan:    healthService.ReportChan(),
	}
	mon.Start()

	if err := microService.Run(); err != nil {
		logger.Fatalf("Run microService failed: %s", err)
	}
}
