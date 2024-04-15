package cameras

import (
	"context"
	"fmt"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/description"
	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/pion/rtp"
	"sync/atomic"
	"time"
)

const maxPackets = 50
const streamTimeout = 15 * time.Second

type cameraStreams struct {
	cam       *rms_cctv.Camera
	Primary   string
	Secondary string
}

func tryStream(ctx context.Context, url string) error {
	cli := gortsplib.Client{}
	uri, err := base.ParseURL(url)
	if err != nil {
		return err
	}

	if err = cli.Start(uri.Scheme, uri.Host); err != nil {
		return err
	}

	defer cli.Close()

	desc, _, err := cli.Describe(uri)
	if err != nil {
		return err
	}

	if err = cli.SetupAll(desc.BaseURL, desc.Medias); err != nil {
		return err
	}

	var packetCnt atomic.Int32
	cli.OnPacketRTPAny(func(medi *description.Media, forma format.Format, pkt *rtp.Packet) {
		packetCnt.Add(1)
	})

	if _, err = cli.Play(nil); err != nil {
		return err
	}

	started := time.Now()
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case now := <-t.C:
			recv := packetCnt.Load()
			if recv >= maxPackets {
				return nil
			}
			if now.Sub(started) >= streamTimeout {
				return fmt.Errorf("low amount RTP packets received: %d < %d", recv, maxPackets)
			}
		}
	}
}
