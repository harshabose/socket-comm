package process

import (
	"context"
	"time"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type StartStreamingHealthSnapshots struct {
	Roomid   types.RoomID  `json:"roomid"`
	Interval time.Duration `json:"interval"`
	AsyncProcess
}

func NewGetHealthSnapshot(roomid types.RoomID) *StartStreamingHealthSnapshots {
	return &StartStreamingHealthSnapshots{
		Roomid: roomid,
	}
}

func (p *StartStreamingHealthSnapshots) Process(ctx context.Context, processor interfaces.Processor, s *state.State) error {
	select {
	case <-ctx.Done():
		return errors.ErrContextCancelled
	default:
		h, ok := processor.(interfaces.CanAddHealthSnapshotStreamer)
		if !ok {
			return errors.ErrInterfaceMisMatch
		}

		return h.AddHealthSnapshotStreamer(p.Roomid, p.Interval, s)
	}
}
