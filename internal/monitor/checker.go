package monitor

import "context"
import "github.com/RacoonMediaServer/rms-packages/pkg/health"

type Checker interface {
	Check(ctx context.Context) []*health.Failure
}
