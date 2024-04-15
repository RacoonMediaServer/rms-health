package cameras

import (
	"context"
	"github.com/RacoonMediaServer/rms-packages/pkg/health"
	"github.com/RacoonMediaServer/rms-packages/pkg/media"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/golang/protobuf/ptypes/empty"
	"go-micro.dev/v4/client"
	"time"
)

const requestTimeout = 20 * time.Second

type Checker struct {
	f           servicemgr.ServiceFactory
	cctvEnabled bool
}

func New(f servicemgr.ServiceFactory, cctvEnabled bool) *Checker {
	return &Checker{
		f:           f,
		cctvEnabled: cctvEnabled,
	}
}

func (c Checker) Check(ctx context.Context) []*health.Failure {
	if !c.cctvEnabled {
		return nil
	}

	cli := c.f.NewCctvCameras()
	resp, err := cli.GetCameras(ctx, &empty.Empty{}, client.WithRequestTimeout(requestTimeout))
	if err != nil {
		return []*health.Failure{makeNotAccessibleFailure(err)}
	}

	var streams []cameraStreams
	var result []*health.Failure

	for _, cam := range resp.Cameras {
		cur := cameraStreams{cam: cam}
		req := rms_cctv.GetLiveUriRequest{
			CameraId:    cam.Id,
			Transport:   media.Transport_RTSP,
			MainProfile: true,
		}
		liveUriResp, err := cli.GetLiveUri(ctx, &req, client.WithRequestTimeout(requestTimeout))
		if err != nil {
			result = append(result, makeLiveStreamIsNotAccessibleFailure(cam, err, true))
		} else {
			cur.Primary = liveUriResp.Uri
		}

		req.MainProfile = false
		liveUriResp, err = cli.GetLiveUri(ctx, &req, client.WithRequestTimeout(requestTimeout))
		if err != nil {
			result = append(result, makeLiveStreamIsNotAccessibleFailure(cam, err, false))
		} else {
			cur.Secondary = liveUriResp.Uri
		}

		streams = append(streams, cur)
	}

	resultChan := make(chan *health.Failure)
	defer close(resultChan)

	workers := 0
	for _, stream := range streams {
		if stream.Primary != "" {
			workers++
			go func(s cameraStreams) {
				if err := tryStream(ctx, s.Primary); err != nil {
					resultChan <- makeLiveStreamPlayFailure(s.cam, err, true)
				} else {
					resultChan <- nil
				}
			}(stream)
		}
		if stream.Secondary != "" {
			workers++
			go func(s cameraStreams) {
				if err := tryStream(ctx, s.Secondary); err != nil {
					resultChan <- makeLiveStreamPlayFailure(s.cam, err, false)
				} else {
					resultChan <- nil
				}
			}(stream)
		}
	}

	for i := 0; i < workers; i++ {
		failure := <-resultChan
		if failure != nil {
			result = append(result, failure)
		}
	}

	return result
}
