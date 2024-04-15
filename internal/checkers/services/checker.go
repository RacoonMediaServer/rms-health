package services

import (
	"context"
	"github.com/RacoonMediaServer/rms-packages/pkg/health"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
)

type Checker struct {
	f           servicemgr.ServiceFactory
	cctvEnabled bool
	required    map[string]struct{}
}

func New(f servicemgr.ServiceFactory, cctvEnabled bool, requiredServices []string) *Checker {
	c := Checker{
		f:           f,
		cctvEnabled: cctvEnabled,
		required:    map[string]struct{}{},
	}

	for _, service := range requiredServices {
		c.required[service] = struct{}{}
	}

	return &c
}

func (c Checker) Check(ctx context.Context) []*health.Failure {
	var result []*health.Failure
	services := c.prepareServiceList()
	for _, srv := range services {
		if srv.CctvService && !c.cctvEnabled {
			continue
		}

		err := srv.Checker(ctx)
		if err != nil {
			result = append(result, c.makeFailure(srv, err))
		}
	}
	return result
}
