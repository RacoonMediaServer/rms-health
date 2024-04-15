package monitor

import (
	"context"
	"github.com/RacoonMediaServer/rms-packages/pkg/health"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"sync"
	"time"
)

type Monitor struct {
	Factory       servicemgr.ServiceFactory
	CheckInterval time.Duration
	Checkers      []Checker
	ReportChan    chan<- *health.Report

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func (m *Monitor) Start() {
	m.ctx, m.cancel = context.WithCancel(context.Background())
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.process()
	}()
}

func (m *Monitor) Stop() {
	m.cancel()
	m.wg.Wait()
}
