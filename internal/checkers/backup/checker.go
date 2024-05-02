package backup

import (
	"context"
	"github.com/RacoonMediaServer/rms-packages/pkg/health"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/golang/protobuf/ptypes/empty"
	"go-micro.dev/v4/client"
	"time"
)

const requestTimeout = 20 * time.Second
const maxBackupExpiration = 24 * 40 * time.Hour

type Checker struct {
	f servicemgr.ServiceFactory
}

func New(f servicemgr.ServiceFactory) *Checker {
	return &Checker{
		f: f,
	}
}

func (c Checker) Check(ctx context.Context) []*health.Failure {
	cli := c.f.NewBackup()
	resp, err := cli.GetBackups(ctx, &empty.Empty{}, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		return []*health.Failure{
			{
				Code:     health.Failure_ServiceNotAccessible,
				Severity: health.Failure_Critical,
				Service:  "rms-backup",
				Text:     err.Error(),
			},
		}
	}
	now := time.Now()
	for i := range resp.Backups {
		backupDate := time.Unix(int64(resp.Backups[i].Date), 0)
		if now.Sub(backupDate) < maxBackupExpiration {
			return nil
		}
	}

	return []*health.Failure{
		{
			Code:     health.Failure_Unknown, // TODO: new failure type
			Severity: health.Failure_Tolerance,
			Service:  "rms-backup",
			Text:     "last backup is expired",
		},
	}
}
