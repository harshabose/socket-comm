package process

import (
	"context"
	"time"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type StartHealthTracking struct {
	RoomID   types.RoomID  `json:"room_id"`
	Interval time.Duration `json:"interval"`
	AsyncProcess
}

func NewStartHealthTracking(roomID types.RoomID) interfaces.CanBeProcessed {
	return &StartHealthTracking{
		RoomID: roomID,
	}
}

func (p *StartHealthTracking) Process(ctx context.Context, processor interfaces.Processor, _ *state.State) error {
	select {
	case <-ctx.Done():
		return errors.ErrContextCancelled
	default:
		t, ok := processor.(interfaces.CanStartHealthTracking)
		if !ok {
			return errors.ErrInterfaceMisMatch
		}

		if err := t.StartHealthTracking(p.RoomID, p.Interval); err != nil {
			return err
		}

		return nil
	}
}
