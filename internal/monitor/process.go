package monitor

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/health"
	"go-micro.dev/v4/logger"
	"time"
)

func (m *Monitor) process() {
	t := time.NewTicker(m.CheckInterval)
	defer t.Stop()
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-t.C:
			m.runCheck()
		}
	}
}

func (m *Monitor) runCheck() {
	logger.Debugf("Perform health checking...")
	ch := make(chan []*health.Failure)
	defer close(ch)
	report := health.Report{
		Status:    health.Report_Ok,
		Timestamp: uint64(time.Now().Unix()),
	}

	for _, checker := range m.Checkers {
		go func(chk Checker) {
			ch <- chk.Check(m.ctx)
		}(checker)
	}
	for range m.Checkers {
		result := <-ch
		report.Failures = append(report.Failures, result...)
	}

	for _, f := range report.Failures {
		if report.Status == health.Report_Ok && f.Severity == health.Failure_Tolerance {
			report.Status = health.Report_Warning
		}
		if f.Severity == health.Failure_Critical {
			report.Status = health.Report_Critical
		}
	}

	m.processReport(&report)
}

func (m *Monitor) processReport(report *health.Report) {
	if report.Status != health.Report_Ok {
		logger.Errorf("Report. Status: %s, Failures: %+v", report.Status.String(), report.Failures)
	}
	if m.ReportChan != nil {
		m.ReportChan <- report
	}
}
