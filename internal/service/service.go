package service

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/health"
	"go-micro.dev/v4"
	"sync"
)

type Service struct {
	pub        micro.Event
	reportChan chan *health.Report

	mu         sync.RWMutex
	lastReport *health.Report
}

func New(pub micro.Event) *Service {
	s := Service{
		pub:        pub,
		reportChan: make(chan *health.Report),
	}

	go s.process()

	return &s
}

func (s *Service) ReportChan() chan<- *health.Report {
	return s.reportChan
}
