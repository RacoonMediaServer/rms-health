package checkers

import (
	"github.com/RacoonMediaServer/rms-health/internal/checkers/cameras"
	"github.com/RacoonMediaServer/rms-health/internal/checkers/containers"
	"github.com/RacoonMediaServer/rms-health/internal/checkers/services"
	"github.com/RacoonMediaServer/rms-health/internal/config"
	"github.com/RacoonMediaServer/rms-health/internal/monitor"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
)

func New(f servicemgr.ServiceFactory, cfg config.Configuration) []monitor.Checker {
	return []monitor.Checker{
		containers.New(f, cfg.RequiredServices),
		cameras.New(f, cfg.Cctv.Enabled),
		services.New(f, cfg.Cctv.Enabled, cfg.RequiredServices),
	}
}
