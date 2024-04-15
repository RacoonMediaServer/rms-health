package containers

import (
	"context"
	"github.com/RacoonMediaServer/rms-packages/pkg/health"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
)

type Checker struct {
	f        servicemgr.ServiceFactory
	required []string
}

func New(f servicemgr.ServiceFactory, requiredServices []string) *Checker {
	return &Checker{
		f:        f,
		required: requiredServices,
	}
}

func (c Checker) Check(ctx context.Context) []*health.Failure {
	runningServices, err := getRunningContainers(ctx)
	if err != nil {
		panic(err)
	}

	var result []*health.Failure
	for _, required := range c.required {
		if _, ok := runningServices[required]; !ok {
			failure := health.Failure{
				Code:     health.Failure_ServiceNotStarted,
				Severity: health.Failure_Critical,
				Service:  required,
				Text:     "Associated docker container not found",
			}
			result = append(result, &failure)
		}
	}
	return result
}
