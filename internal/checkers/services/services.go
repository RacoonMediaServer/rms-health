package services

import (
	"context"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	rms_torrent "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-torrent"
	"github.com/golang/protobuf/ptypes/empty"
	"go-micro.dev/v4/client"
	"time"
)

const requestTimeout = 10 * time.Second

type serviceInfo struct {
	ContainerName string
	ServiceName   string
	CctvService   bool
	Checker       func(ctx context.Context) error
}

func (c Checker) checkCctv(ctx context.Context) error {
	_, err := c.f.NewCctvCameras().GetCameras(ctx, &empty.Empty{}, client.WithRequestTimeout(requestTimeout))
	return err
}

func (c Checker) checkBotClient(ctx context.Context) error {
	_, err := c.f.NewBotClient().GetIdentificationCode(ctx, &empty.Empty{}, client.WithRequestTimeout(requestTimeout))
	return err
}

func (c Checker) checkTranscoder(ctx context.Context) error {
	_, err := c.f.NewTranscoderProfiles().GetProfiles(ctx, &empty.Empty{}, client.WithRequestTimeout(requestTimeout))
	return err
}

func (c Checker) checkNotes(ctx context.Context) error {
	_, err := c.f.NewNotes().GetSettings(ctx, &empty.Empty{})
	return err
}

func (c Checker) checkLibrary(ctx context.Context) error {
	_, err := c.f.NewLibrary().GetMovies(ctx, &rms_library.GetMoviesRequest{}, client.WithRequestTimeout(requestTimeout))
	return err
}

func (c Checker) checkTorrent(ctx context.Context) error {
	_, err := c.f.NewTorrent().GetTorrents(ctx, &rms_torrent.GetTorrentsRequest{IncludeDoneTorrents: false}, client.WithRequestTimeout(requestTimeout))
	return err
}

func (c Checker) checkNotifier(ctx context.Context) error {
	_, err := c.f.NewNotifier().GetSettings(ctx, &empty.Empty{}, client.WithRequestTimeout(requestTimeout))
	return err
}

func (c Checker) checkBackup(ctx context.Context) error {
	_, err := c.f.NewBackup().GetBackupSettings(ctx, &empty.Empty{}, client.WithRequestTimeout(requestTimeout))
	return err
}

func (c Checker) prepareServiceList() []serviceInfo {
	return []serviceInfo{
		{
			ContainerName: "cctv",
			ServiceName:   "rms-cctv",
			CctvService:   true,
			Checker:       c.checkCctv,
		},
		{
			ContainerName: "bot-client",
			ServiceName:   "rms-bot-client",
			CctvService:   false,
			Checker:       c.checkBotClient,
		},
		{
			ContainerName: "transcoder",
			ServiceName:   "rms-transcoder",
			CctvService:   false,
			Checker:       c.checkTranscoder,
		},
		{
			ContainerName: "notes",
			ServiceName:   "rms-notes",
			CctvService:   false,
			Checker:       c.checkNotes,
		},
		{
			ContainerName: "library",
			ServiceName:   "rms-library",
			CctvService:   false,
			Checker:       c.checkLibrary,
		},
		{
			ContainerName: "torrent",
			ServiceName:   "rms-torrent",
			CctvService:   false,
			Checker:       c.checkTorrent,
		},
		{
			ContainerName: "notifier",
			ServiceName:   "rms-notifier",
			CctvService:   false,
			Checker:       c.checkNotifier,
		},
		{
			ContainerName: "backup",
			ServiceName:   "rms-backup",
			CctvService:   false,
			Checker:       c.checkBackup,
		},
	}
}
