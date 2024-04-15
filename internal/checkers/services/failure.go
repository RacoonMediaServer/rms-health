package services

import "github.com/RacoonMediaServer/rms-packages/pkg/health"

func (c Checker) makeFailure(service serviceInfo, err error) *health.Failure {
	severity := health.Failure_Tolerance
	if _, ok := c.required[service.ContainerName]; ok {
		severity = health.Failure_Critical
	}
	return &health.Failure{
		Code:     health.Failure_ServiceNotAccessible,
		Severity: severity,
		Service:  service.ServiceName,
		Text:     err.Error(),
	}
}
