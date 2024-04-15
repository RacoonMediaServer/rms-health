package service

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"github.com/RacoonMediaServer/rms-packages/pkg/health"
)

func getCriticalFailure(report *health.Report) *health.Failure {
	for _, f := range report.Failures {
		if f.Severity == health.Failure_Critical {
			return f
		}
	}
	return &health.Failure{}
}

func decodeMalfunction(report *health.Report) (events.Malfunction_System, events.Malfunction_Code) {
	f := getCriticalFailure(report)
	switch f.Code {
	case health.Failure_ServiceNotStarted:
		fallthrough
	case health.Failure_ServiceNotAccessible:
		return events.Malfunction_Services, events.Malfunction_CannotAccess
	case health.Failure_LiveStreamIsNotAccessible:
		fallthrough
	case health.Failure_LiveStreamPlayFailure:
		return events.Malfunction_Cameras, events.Malfunction_CannotAccess
	case health.Failure_RecordingUnexpectedStopped:
		return events.Malfunction_Archive, events.Malfunction_TaskIsHung
	case health.Failure_RecordingPlaybackFailure:
		return events.Malfunction_Archive, events.Malfunction_CannotAccess
	}

	return events.Malfunction_Services, events.Malfunction_Unknown
}
