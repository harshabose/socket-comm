package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type StopStreamingHealthSnapshots struct {
	Roomid types.RoomID `json:"roomid"`
	AsyncProcess
}

func NewStopStreamingHealthSnapshots(roomid types.RoomID) StopStreamingHealthSnapshots {
	return StopStreamingHealthSnapshots{
		Roomid: roomid,
	}
}

func (p *StopStreamingHealthSnapshots) Process(ctx context.Context, processor interfaces.Processor, s *state.State) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		h, ok := processor.(interfaces.CanRemoveHealthSnapshotStreamer)
		if !ok {
			return errors.ErrInterfaceMisMatch
		}

		return h.RemoveHealthSnapshotStreamer(p.Roomid, s)
	}
}
