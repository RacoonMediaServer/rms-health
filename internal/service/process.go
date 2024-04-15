package service

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"github.com/RacoonMediaServer/rms-packages/pkg/health"
	"go-micro.dev/v4/logger"
	"time"
)

const notifyTimeout = 20 * time.Second

func (s *Service) process() {
	for {
		report := <-s.reportChan
		s.mu.Lock()
		notify := s.lastReport != nil && s.lastReport.Status != report.Status
		s.lastReport = report
		s.mu.Unlock()

		if notify {
			s.notify(report)
		}
	}
}

func (s *Service) notify(report *health.Report) {
	if report.Status == health.Report_Critical {
		event := events.Malfunction{
			Sender:     "rms-health",
			Timestamp:  int64(report.Timestamp),
			Error:      getCriticalFailure(report).Code.String(),
			StackTrace: fmt.Sprintf("%+v", report.Failures),
		}
		event.System, event.Code = decodeMalfunction(report)

		ctx, cancel := context.WithTimeout(context.Background(), notifyTimeout)
		defer cancel()

		if err := s.pub.Publish(ctx, &event); err != nil {
			logger.Errorf("Notify failed: %s", err)
		}
	}
}
