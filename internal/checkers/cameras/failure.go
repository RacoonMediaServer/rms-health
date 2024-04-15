package cameras

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-packages/pkg/health"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
)

func makeNotAccessibleFailure(err error) *health.Failure {
	return &health.Failure{
		Code:     health.Failure_ServiceNotAccessible,
		Severity: health.Failure_Critical,
		Service:  "rms-cctv",
		Text:     err.Error(),
	}
}

func makeLiveStreamIsNotAccessibleFailure(cam *rms_cctv.Camera, err error, isMain bool) *health.Failure {
	id := fmt.Sprintf("%d", cam.Id)
	streamType := "primary"
	if !isMain {
		streamType = "secondary"
	}
	return &health.Failure{
		Code:        health.Failure_LiveStreamIsNotAccessible,
		Severity:    health.Failure_Critical,
		Service:     "rms-cctv",
		Text:        err.Error(),
		SubjectId:   &id,
		SubjectName: &cam.Name,
		Details:     map[string]string{"streamType": streamType},
	}
}

func makeLiveStreamPlayFailure(cam *rms_cctv.Camera, err error, isMain bool) *health.Failure {
	id := fmt.Sprintf("%d", cam.Id)
	streamType := "primary"
	if !isMain {
		streamType = "secondary"
	}
	return &health.Failure{
		Code:        health.Failure_LiveStreamPlayFailure,
		Severity:    health.Failure_Critical,
		Service:     "rms-cctv",
		Text:        err.Error(),
		SubjectId:   &id,
		SubjectName: &cam.Name,
		Details:     map[string]string{"streamType": streamType},
	}
}
